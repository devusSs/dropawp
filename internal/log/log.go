package log

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Setup(dir string) error {
	if dir == "" {
		return errors.New("log directory cannot be empty")
	}

	err := createLoggers(dir)
	if err != nil {
		return fmt.Errorf("failed to create loggers: %w", err)
	}

	return nil
}

func Debug(msg string, args ...any) {
	if consoleLogger == nil || fileLogger == nil {
		fmt.Fprintln(os.Stderr, "Loggers are not initialized. Please call Setup()")
		return
	}

	consoleLogger.Debug(msg, args...)
	fileLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if consoleLogger == nil || fileLogger == nil {
		fmt.Fprintln(os.Stderr, "Loggers are not initialized. Please call Setup()")
		return
	}

	consoleLogger.Info(msg, args...)
	fileLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if consoleLogger == nil || fileLogger == nil {
		fmt.Fprintln(os.Stderr, "Loggers are not initialized. Please call Setup()")
		return
	}

	consoleLogger.Warn(msg, args...)
	fileLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if consoleLogger == nil || fileLogger == nil {
		fmt.Fprintln(os.Stderr, "Loggers are not initialized. Please call Setup()")
		return
	}

	consoleLogger.Error(msg, args...)
	fileLogger.Error(msg, args...)
}

var (
	consoleLogger *slog.Logger
	fileLogger    *slog.Logger
)

func createLoggers(dir string) error {
	fileWriter, err := createCurrentLogFile(dir)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	consoleLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: getLogLevelFromEnv(),
	}))

	fileLogger = slog.New(slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
		Level: getLogLevelFromEnv(),
	}))

	return nil
}

func createCurrentLogFile(dir string) (io.Writer, error) {
	date := time.Now().Format("2006-01-02_15-04-05")
	name := fmt.Sprintf("dropawp_%s.log", date)

	err := createDirs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	path := filepath.Join(dir, name)

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", path, err)
	}

	return f, nil
}

func createDirs(dir string) error {
	if dir == "" {
		return errors.New("directory path cannot be empty")
	}

	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	return nil
}

func getLogLevelFromEnv() slog.Level {
	l := strings.ToLower(os.Getenv("DROPAWP_LOG_LEVEL"))

	switch l {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
