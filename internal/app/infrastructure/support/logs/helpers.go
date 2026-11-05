package logs

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var (
	apiLogPath = filepath.Join("logs", "api.log")
	dbLogPath  = filepath.Join("logs", "db.log")
)

// format converts any value to a string for logging. Primitives are written
// as-is; structs, maps and slices are serialized as JSON.
func format(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case nil:
		return "null"
	case []byte:
		return string(val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return "null"
		}
		return string(b)
	}
}

// prependLine reads the existing file content, then writes the new entry
// at the top of the file. Each write opens and closes the file so no handle
// is held between reads.
func prependLine(path, newLine string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Warn("could not create logs directory", "path", path, "error", err)
		return err
	}

	content, _ := os.ReadFile(path)

	header := time.Now().UTC().Format(time.RFC3339Nano) + " | " + newLine

	var combined string
	if len(content) > 0 {
		combined = header + "\n" + string(content)
	} else {
		combined = header
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(combined)
	file.Close()
	return err
}
