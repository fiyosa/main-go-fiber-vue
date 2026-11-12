package database

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-fiber-vue/internal/app/infrastructure/support/logs"
)

// Open connects to PostgreSQL. When logSQL is true every executed statement is
// written to logs/db.log (newest at the top); otherwise only slow queries and
// errors are reported.
func Open(dsn string, logSQL bool) (*gorm.DB, error) {
	gormLogger := newLogger(logSQL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

// newLogger returns the GORM logger. With DB_LOG_SQL enabled, all statements
// go to logs.LogDB so they land in logs/db.log with newest at the top.
func newLogger(logSQL bool) logger.Interface {
	if !logSQL {
		return logger.Default.LogMode(logger.Warn)
	}
	return &sqlLogger{logMode: logger.Info}
}

type sqlLogger struct {
	logMode logger.LogLevel
}

func (s *sqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &sqlLogger{logMode: level}
}

func (s *sqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logs.LogDB(msg)
}

func (s *sqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logs.LogDB(msg)
}

func (s *sqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logs.LogDB(msg)
}

func (s *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if s.logMode <= logger.Silent {
		return
	}
	sql, _ := fc()
	if sql != "" {
		logs.LogDB(sql)
	}
}
