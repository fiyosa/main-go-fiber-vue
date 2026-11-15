package database

import (
	"context"
	"fmt"
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

// slowQueryThreshold is the duration above which a statement is written to
// logs/db.log while DB_LOG_SQL is disabled. Slow statements and query errors go
// to that file so the terminal stays clean without losing diagnostics.
const slowQueryThreshold = 200 * time.Millisecond

// newLogger returns the GORM logger. With DB_LOG_SQL enabled every statement
// goes to logs.LogDB so they land in logs/db.log with newest at the top;
// otherwise only slow queries and errors are reported, and they go to the same
// file so GORM never writes to stdout.
func newLogger(logSQL bool) logger.Interface {
	if logSQL {
		return &sqlLogger{logMode: logger.Info}
	}
	return &sqlLogger{logMode: logger.Warn, slowThreshold: slowQueryThreshold}
}

type sqlLogger struct {
	logMode       logger.LogLevel
	slowThreshold time.Duration
}

func (s *sqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *s
	clone.logMode = level
	return &clone
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

	elapsed := time.Since(begin)

	switch {
	case err != nil && s.logMode >= logger.Error:
		sql, rows := fc()
		logs.LogDB(formatTrace("error", elapsed, rows, err, sql))
	case s.slowThreshold > 0 && elapsed > s.slowThreshold && s.logMode >= logger.Warn:
		sql, rows := fc()
		logs.LogDB(formatTrace("slow", elapsed, rows, nil, sql))
	case s.logMode >= logger.Info:
		if sql, _ := fc(); sql != "" {
			logs.LogDB(sql)
		}
	}
}

// formatTrace builds a single log line so slow queries and query errors can be
// filtered with grep, e.g. grep "\[slow\]" logs/db.log.
func formatTrace(tag string, elapsed time.Duration, rows int64, err error, sql string) string {
	line := fmt.Sprintf("[%s] %.3fms rows:%d | %s", tag, float64(elapsed.Microseconds())/1000, rows, sql)
	if err != nil {
		return line + " | error: " + err.Error()
	}
	return line
}
