package logs

// LogDB writes entry to ./logs/db.log with the newest line at the top.
func LogDB(entry any) error {
	return prependLine(dbLogPath, format(entry))
}
