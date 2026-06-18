package utils

import "time"

// Now returns the current time formatted as "YYYY-MM-DD HH:MM:SS".
// Usage: utils.Now()
// Output: "2026-06-08 10:30:00"
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatByDate formats a time.Time using Go's reference layout.
// Layout: "2006-01-02 15:04:05", "2006-01-02", "15:04:05", etc.
// Usage: utils.FormatByDate(time.Now(), "2006-01-02 15:04:05")
// Output: "2026-06-08 10:30:00"
func FormatByDate(t time.Time, format string) string {
	return t.Format(format)
}

// FormatByStr parses a date string from one layout and reformats it to another.
// Usage: utils.FormatByStr("2024-01-15", "2006-01-02", "02/01/2006")
// Output: "15/01/2024", nil
func FormatByStr(dateStr, fromFormat, toFormat string) (string, error) {
	t, err := time.Parse(fromFormat, dateStr)
	if err != nil {
		return "", err
	}
	return t.Format(toFormat), nil
}
