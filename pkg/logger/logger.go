package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type LoggerType string

func NewLogger(ctx *context.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create new logger: %v", err)
	}
	zLog := logger.Sugar()

	// Create context with logger
	var (
		loggerName LoggerType
	)
	loggerName = "ZapLogger"
	*ctx = context.WithValue(*ctx, loggerName, zLog)
	return nil
}
