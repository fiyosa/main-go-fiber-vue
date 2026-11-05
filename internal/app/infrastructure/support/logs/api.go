package logs

// Log writes entry to ./logs/api.log with the newest line at the top.
func Log(entry any) error {
	return prependLine(apiLogPath, format(entry))
}
