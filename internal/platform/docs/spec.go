package docs

import (
	"encoding/json"
	"sort"
	"strings"
)

// RouteInfo represents a registered route path and HTTP method.
type RouteInfo struct {
	Method string
	Path   string
}

// IsAPIRoute returns true for API routes under /api/ but excludes MCP and docs.
func IsAPIRoute(path, mcpPath string) bool {
	if path == mcpPath {
		return false
	}
	if strings.HasPrefix(path, "/api/docs") || strings.HasPrefix(path, "/docs") {
		return false
	}
	return path == "/api" || strings.HasPrefix(path, "/api/")
}

// BuildSpec builds an OpenAPI 3.1 specification from discovered API routes,
// merged with base defaults and optional manual overrides.
func BuildSpec(routes []RouteInfo, overrideJSON []byte) (map[string]any, error) {
	// Base OpenAPI 3.1 document defaults
	spec := map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "app-api",
			"version":     "0.1.0",
			"description": "API Reference for app-api endpoints",
		},
		"servers": []any{
			map[string]any{
				"url":         "",
				"description": "Relative server root",
			},
		},
		"paths": map[string]any{},
		"components": map[string]any{
			"schemas": map[string]any{
				"SuccessEnvelope": map[string]any{
					"type":     "object",
					"required": []any{"success", "message"},
					"properties": map[string]any{
						"success": map[string]any{"type": "boolean", "example": true},
						"data":    map[string]any{"type": "object"},
						"message": map[string]any{"type": "string"},
					},
				},
				"ErrorEnvelope": map[string]any{
					"type":     "object",
					"required": []any{"success", "message"},
					"properties": map[string]any{
						"success": map[string]any{"type": "boolean", "example": false},
						"message": map[string]any{"type": "string"},
					},
				},
				"ValidationErrorEnvelope": map[string]any{
					"type":     "object",
					"required": []any{"success", "message", "errors"},
					"properties": map[string]any{
						"success": map[string]any{"type": "boolean", "example": false},
						"message": map[string]any{"type": "string"},
						"errors": map[string]any{
							"type": "object",
							"additionalProperties": map[string]any{
								"type": "array",
								"items": map[string]any{"type": "string"},
							},
						},
					},
				},
				"LoginRequest": map[string]any{
					"type":     "object",
					"required": []any{"email", "password"},
					"properties": map[string]any{
						"email":    map[string]any{"type": "string", "format": "email", "example": "user@example.com"},
						"password": map[string]any{"type": "string", "format": "password", "example": "secret123"},
					},
				},
				"LoginResponse": map[string]any{
					"type":     "object",
					"required": []any{"success", "message"},
					"properties": map[string]any{
						"success": map[string]any{"type": "boolean", "example": true},
						"message": map[string]any{"type": "string", "example": "Logged in successfully"},
					},
				},
			},
			"securitySchemes": map[string]any{
				"cookieAuth": map[string]any{
					"type": "apiKey",
					"in":   "cookie",
					"name": "access_token",
				},
			},
		},
	}

	// Sort routes for deterministic output
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	paths := spec["paths"].(map[string]any)

	// Populate paths from discovered routes
	for _, r := range routes {
		pathObj, ok := paths[r.Path].(map[string]any)
		if !ok {
			pathObj = map[string]any{}
			paths[r.Path] = pathObj
		}

		methodLower := strings.ToLower(r.Method)
		// Skip unsupported HTTP methods in spec if any
		if methodLower == "" {
			continue
		}

		// Default operation stub or rich metadata for known endpoints
		op := defaultOperation(r.Path, methodLower)
		pathObj[methodLower] = op
	}

	// Apply manual override if provided
	if len(overrideJSON) > 0 {
		var override map[string]any
		if err := json.Unmarshal(overrideJSON, &override); err != nil {
			return nil, err
		}
		mergeMap(spec, override)
	}

	return spec, nil
}

func defaultOperation(path, method string) map[string]any {
	// Known endpoints rich definitions
	if path == "/api/v1/auth/login" && method == "post" {
		return map[string]any{
			"tags":        []any{"Auth"},
			"summary":     "Login",
			"description": "Authenticates user credentials. On success, sets HTTP-only cookies (`access_token` and `refresh_token`).",
			"requestBody": map[string]any{
				"required": true,
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": map[string]any{"$ref": "#/components/schemas/LoginRequest"},
					},
				},
			},
			"responses": map[string]any{
				"200": map[string]any{
					"description": "Login successful. Auth cookies issued.",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema":  map[string]any{"$ref": "#/components/schemas/LoginResponse"},
							"example": map[string]any{"success": true, "message": "Logged in successfully"},
						},
					},
				},
				"401": map[string]any{
					"description": "Invalid credentials",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema":  map[string]any{"$ref": "#/components/schemas/ErrorEnvelope"},
							"example": map[string]any{"success": false, "message": "Invalid email or password"},
						},
					},
				},
				"422": map[string]any{
					"description": "Validation error",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema":  map[string]any{"$ref": "#/components/schemas/ValidationErrorEnvelope"},
							"example": map[string]any{"success": false, "message": "The given data was invalid."},
						},
					},
				},
			},
		}
	}
	if path == "/api/v1/auth/logout" && method == "post" {
		return map[string]any{
			"tags":        []any{"Auth"},
			"summary":     "Logout",
			"description": "Revokes the current session and clears the HTTP-only `access_token` and `refresh_token` cookies.",
			"responses": map[string]any{
				"200": map[string]any{
					"description": "Logout successful. Auth cookies cleared.",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema":  map[string]any{"$ref": "#/components/schemas/LoginResponse"},
							"example": map[string]any{"success": true, "message": "Logged out successfully"},
						},
					},
				},
			},
		}
	}
	if (path == "/api/v1/healthz" || path == "/api/v1/readyz") && method == "get" {
		tag := "Health"
		summary := "Liveness probe"
		if path == "/api/v1/readyz" {
			summary = "Readiness probe"
		}
		return map[string]any{
			"tags":        []any{tag},
			"summary":     summary,
			"description": "Probe endpoint.",
			"responses": map[string]any{
				"200": map[string]any{
					"description": "Success",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{"$ref": "#/components/schemas/SuccessEnvelope"},
						},
					},
				},
			},
		}
	}

	// Generic stub for any other discovered route
	return map[string]any{
		"tags":    []any{"Default"},
		"summary": strings.ToUpper(method) + " " + path,
		"responses": map[string]any{
			"200": map[string]any{
				"description": "Successful response",
			},
		},
	}
}

// mergeMap recursively merges override into base.
// If override value is nil, the key is deleted from base.
// Maps are merged recursively, while slices and scalars are replaced.
func mergeMap(base, override map[string]any) {
	for k, ov := range override {
		if ov == nil {
			delete(base, k)
			continue
		}

		bv, exists := base[k]
		if !exists {
			base[k] = ov
			continue
		}

		bMap, bIsMap := bv.(map[string]any)
		oMap, oIsMap := ov.(map[string]any)

		if bIsMap && oIsMap {
			mergeMap(bMap, oMap)
		} else {
			// Replace scalar, slice, or mixed types
			base[k] = ov
		}
	}
}
