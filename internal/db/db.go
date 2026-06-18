package db

import (
	"fmt"
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/lib"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var RUN *gorm.DB

type dbWriter struct{}

func (dbWriter) Printf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	parts := strings.SplitN(msg, "\n", 2)
	if len(parts) > 1 {
		msg = parts[1]
	}
	if idx := strings.LastIndex(msg, "] "); idx != -1 {
		msg = strings.TrimSpace(msg[idx+2:])
	}
	lib.Log.Info(strings.ReplaceAll(msg, `"`, "'"), "db")
}

func Init() {
	var err error
	RUN, err = gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{
		Logger: logger.New(
			dbWriter{},
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	lib.Log.Info("Database connected", "db")
}
