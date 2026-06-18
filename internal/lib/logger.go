package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"go-fiber-svelte/internal/config"

	"github.com/rs/zerolog"
)

var Log logLib

type logFile struct {
	writer *logWriter
	logger zerolog.Logger
}

type logLib struct {
	mu    sync.Mutex
	files map[string]*logFile
}

type logEntry struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type logWriter struct {
	file     *os.File
	mu       sync.Mutex
	logName  string
	filePath string
}

func (l *logLib) getOrCreate(name string) *zerolog.Logger {
	if !config.APP_LOG {
		logger := zerolog.Nop()
		return &logger
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.files == nil {
		l.files = make(map[string]*logFile)
	}
	if f, ok := l.files[name]; ok {
		return &f.logger
	}

	logDir := "./logs"
	os.MkdirAll(logDir, 0755)
	file, err := os.OpenFile(logDir+"/"+name+"_"+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		l.files[name] = &logFile{logger: logger}
		return &logger
	}
	writer := &logWriter{file: file, logName: name, filePath: logDir + "/" + name + "_" + time.Now().Format("2006-01-02") + ".log"}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	l.files[name] = &logFile{writer: writer, logger: logger}
	return &logger
}

func (l *logLib) CloseFile(filePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for name, f := range l.files {
		if f.writer != nil && f.writer.filePath == filePath {
			f.writer.file.Close()
			delete(l.files, name)
			return
		}
	}
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var raw map[string]any
	if err := json.Unmarshal(p, &raw); err != nil {
		return 0, err
	}

	t, _ := raw["time"].(string)
	parsed, _ := time.Parse(time.RFC3339, t)

	entry := logEntry{
		Level:   fmt.Sprintf("%s", raw["level"]),
		Time:    parsed.Format("2006-01-02 15:04:05"),
		Message: fmt.Sprintf("%s", raw["message"]),
	}

	data, _ := json.Marshal(entry)
	w.file.Write(data)
	w.file.WriteString("\n")
	return len(p), nil
}

func (l *logLib) Info(msg string, name ...string) {
	logName := "fiber"
	if len(name) > 0 {
		logName = name[0]
	}
	l.getOrCreate(logName).Info().Msg(msg)
}

func (l *logLib) Error(msg string, err error, name ...string) {
	logName := "fiber"
	if len(name) > 0 {
		logName = name[0]
	}
	l.getOrCreate(logName).Error().Err(err).Msg(msg)
}
