package logger

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

func NewGormLogger() logger.Interface {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[INFO] "+msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[WARN] "+msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[ERROR] "+msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	file := utils.FileWithLineNum()

	if err != nil {
		log.Printf("[ERROR] %s | %v | %s | %d rows | %s", file, err, elapsed, rows, sql)
	} else if elapsed > l.SlowThreshold {
		log.Printf("[SLOW SQL] %s | %s | %d rows | %s", file, elapsed, rows, sql)
	} else {
		log.Printf("[QUERY] %s | %s | %d rows | %s", file, elapsed, rows, sql)
	}
}
