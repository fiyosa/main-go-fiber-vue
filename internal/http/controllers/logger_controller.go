package controllers

import (
	"go-fiber-svelte/internal/http/repositories/logger_repository"

	"github.com/gofiber/fiber/v2"
)

func LoggerList(c *fiber.Ctx) error {
	return logger_repository.LogListRepository(c)
}

func LoggerDetail(c *fiber.Ctx) error {
	return logger_repository.LogDetailRepository(c)
}

func LoggerDownload(c *fiber.Ctx) error {
	return logger_repository.LogDownloadRepository(c)
}

func LoggerDelete(c *fiber.Ctx) error {
	return logger_repository.LogDeleteRepository(c)
}

type logPathItem struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Summary string `json:"summary"`
}

func LoggerOpenAPIPaths() map[string]any {
	return map[string]any{
		"/api/log": map[string]any{
			"get": map[string]any{
				"summary":     "List log files",
				"description": "Returns list of available log files with size",
				"tags":        []string{"Log"},
				"responses": map[string]any{
					"200": map[string]any{"description": "List of log files"},
				},
			},
		},
		"/api/log/{filename}": map[string]any{
			"get": map[string]any{
				"summary":     "Get log detail",
				"description": "Returns log entries from a specific file, with optional level filter, search, and pagination",
				"tags":        []string{"Log"},
				"parameters": []map[string]any{
					{"name": "filename", "in": "path", "required": true, "schema": map[string]any{"type": "string"}},
					{"name": "levels", "in": "query", "schema": map[string]any{"type": "string"}, "description": "Comma-separated level filter"},
					{"name": "search", "in": "query", "schema": map[string]any{"type": "string"}, "description": "Search text"},
					{"name": "page", "in": "query", "schema": map[string]any{"type": "integer"}, "description": "Page number"},
					{"name": "limit", "in": "query", "schema": map[string]any{"type": "integer"}, "description": "Items per page"},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Log entries"},
				},
			},
			"delete": map[string]any{
				"summary":     "Delete log file",
				"description": "Deletes a log file",
				"tags":        []string{"Log"},
				"parameters": []map[string]any{
					{"name": "filename", "in": "path", "required": true, "schema": map[string]any{"type": "string"}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Log file deleted"},
				},
			},
		},
		"/api/log/{filename}/download": map[string]any{
			"get": map[string]any{
				"summary":     "Download log file",
				"description": "Downloads a log file",
				"tags":        []string{"Log"},
				"parameters": []map[string]any{
					{"name": "filename", "in": "path", "required": true, "schema": map[string]any{"type": "string"}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Log file download"},
				},
			},
		},
	}
}
