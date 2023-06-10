package logging

import (
	"context"
	"log"
	"os"
	"path/filepath"
)

type loggerKeyType int

const LoggerKey loggerKeyType = iota

/**
 * InitLogger initializes the logger which prints the logs in a local file and adds it to the context.
 */
func InitLogger(ctx context.Context, logFilePath string) context.Context {
	logFileDir := filepath.Dir(logFilePath)

	err := os.MkdirAll(logFileDir, 0755)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", log.LstdFlags)

	ctx = context.WithValue(ctx, LoggerKey, logger)

	return ctx
}

/**
 * GetLogger returns the logger from the context.
 */
func GetLogger(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*log.Logger); ok {
		return logger
	}
	return nil
}
