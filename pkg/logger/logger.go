package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type LoggerType string

func NewLogger(ctx *context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to create new logger: %v", err)
	}
	zLog := logger.Sugar()

	// Create context with logger
	var loggerName LoggerType = "ZapLogger"
	*ctx = context.WithValue(*ctx, loggerName, zLog)
	return nil
}

func GetLogger(ctx *context.Context) *zap.SugaredLogger {
	var loggerName LoggerType = "ZapLogger"
	loggerName = "ZapLogger"
	zLog := (*ctx).Value(loggerName).(*zap.SugaredLogger)
	return zLog
}
