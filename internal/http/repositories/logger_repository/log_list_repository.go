package logger_repository

import (
	"os"
	"sort"
	"strings"

	"go-fiber-svelte/internal/helper"

	"github.com/gofiber/fiber/v2"
)

type LogFileItem struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func LogListRepository(c *fiber.Ctx) error {
	dir := "./logs"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return c.JSON(helper.Res.SuccessData([]LogFileItem{}, "Log files retrieved successfully"))
	}

	var files []LogFileItem
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		info, _ := e.Info()
		size := int64(0)
		if info != nil {
			size = info.Size()
		}
		files = append(files, LogFileItem{Name: e.Name(), Size: size})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})

	return c.JSON(helper.Res.SuccessData(files, "Log files retrieved successfully"))
}
