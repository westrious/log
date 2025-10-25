package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
)

var (
	logger        *slog.Logger
	noticeKVStore []any

	levelFatal    = slog.LevelError + 999
	levelNoticeKV = slog.LevelInfo + 1
)

const size = 64 << 20 // 64 MB

type Options struct {
	ProjectName string       `json:"project_name"`
	LogDir      string       `json:"log_dir"`
	Level       slog.Leveler `json:"level"`
	AddSource   bool         `json:"add_source"`
}

func Init(options Options) {
	if options.ProjectName == "" {
		options.ProjectName = "app"
	}
	if options.LogDir == "" {
		options.LogDir = "./logs"
	}
	if options.Level == nil {
		options.Level = slog.LevelInfo
	}

	if err := os.MkdirAll(options.LogDir, 0766); err != nil {
		panic(err)
	}
	logPath := path.Join(options.LogDir, fmt.Sprintf("%s.log", options.ProjectName))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	logger = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	}))

	noticeKVStore = make([]any, 0)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Debugf(format string, a ...any) {
	logger.Debug(fmt.Sprintf(format, a...))
}

func Infof(format string, a ...any) {
	logger.Info(fmt.Sprintf(format, a...))
}

func Warnf(format string, a ...any) {
	logger.Warn(fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...any) {
	logger.Error(fmt.Sprintf(format, a...))
}

func Fatal(msg string, err any) {
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	logger.Log(context.Background(), levelFatal, msg, slog.Any("error", err), slog.String("stacktrace", string(buf)))
	os.Exit(1)
}

func Flush() {
	if len(noticeKVStore) == 0 {
		return
	}
	logger.Log(context.Background(), levelNoticeKV, "NoticeKV", noticeKVStore...)
}

func PushNotice(key string, value any) {
	noticeKVStore = append(noticeKVStore, key, value)
}
