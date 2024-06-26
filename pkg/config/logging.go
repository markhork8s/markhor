package config

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var LoggerLevels = map[string]slog.Level{
	"info":    slog.LevelInfo,
	"debug":   slog.LevelDebug,
	"warning": slog.LevelWarn,
	"error":   slog.LevelError,
}

func SetupLogging(config Config) error {

	files := make([]*os.File, 0)

	if config.Logging.LogToStdout {
		files = append(files, os.Stdout)
	}

	for _, fname := range config.Logging.AdditionalLogFiles {
		//#nosec G304 -- This is the intended functionality (gosec)
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		files = append(files, f)
	}

	writers := make([]io.Writer, len(files))
	for i, file := range files {
		writers[i] = file
	}

	writer := io.MultiWriter(writers...)

	level, ok := LoggerLevels[config.Logging.Level]
	if !ok {
		return errors.New(fmt.Sprint("Invalid log level specified:", config.Logging.Level))
	}
	isDebugModeActive := level == slog.LevelDebug
	opts := &slog.HandlerOptions{
		AddSource: isDebugModeActive,
		Level:     level,
	}
	var handler slog.Handler
	switch config.Logging.Style {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	case "json":
		handler = slog.NewJSONHandler(writer, opts)
	default:
		return errors.New(fmt.Sprint("Invalid log style specified:", config.Logging.Style))
	}

	slog.SetDefault(slog.New(handler))
	return nil
}
