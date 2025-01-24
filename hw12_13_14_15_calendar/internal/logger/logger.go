package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	file   *os.File
	logger *slog.Logger
}

func New(level string, fileName string) *Logger {
	// setup:
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var programLevel = new(slog.LevelVar)
	if err := programLevel.UnmarshalText([]byte(level)); err != nil {
		log.Printf("Error parsing level: %v. Acceptable values: DEBUG, INFO, WARN, ERROR. Will use INFO\n", err)
		programLevel.Set(slog.LevelInfo)
	}
	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		Level:       programLevel,
		ReplaceAttr: nil,
	}
	writer := io.MultiWriter(file, os.Stderr)
	slogger := slog.New(slog.NewTextHandler(writer, logConfig))
	slog.SetDefault(slogger)

	return &Logger{file: file, logger: slogger}
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Close() error {
	return l.file.Close()
}
