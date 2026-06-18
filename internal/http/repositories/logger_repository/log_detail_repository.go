package logger_repository

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"go-fiber-svelte/internal/helper"

	"github.com/gofiber/fiber/v2"
)

type LogEntry struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type LogDetailResponse struct {
	Name  string     `json:"name"`
	Total int        `json:"total"`
	Logs  []LogEntry `json:"logs"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

func LogDetailRepository(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" || strings.Contains(filename, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(helper.Res.Error("Invalid filename", nil))
	}

	file, err := os.Open("./logs/" + filename)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(helper.Res.Error("Log file not found", nil))
	}
	defer file.Close()

	levels := c.Query("levels")
	search := strings.ToLower(c.Query("search"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 500 {
		limit = 50
	}

	var levelSet map[string]bool
	if levels != "" {
		levelSet = make(map[string]bool)
		for _, l := range strings.Split(levels, ",") {
			levelSet[strings.TrimSpace(l)] = true
		}
	}

	allEntries := make([]LogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if levelSet != nil && !levelSet[entry.Level] {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(entry.Message), search) && !strings.Contains(strings.ToLower(entry.Level), search) {
			continue
		}
		allEntries = append(allEntries, entry)
	}

	if err := scanner.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.Res.Error("Failed to read log file", nil))
	}

	for i, j := 0, len(allEntries)-1; i < j; i, j = i+1, j-1 {
		allEntries[i], allEntries[j] = allEntries[j], allEntries[i]
	}

	total := len(allEntries)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return c.JSON(helper.Res.Paginate(allEntries[start:end], &helper.Meta{
		Total: total,
		Page:  page,
		Limit: limit,
	}, "Log detail retrieved successfully"))
}
