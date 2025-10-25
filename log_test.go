package log

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

// helper to initialize logger to a temporary file
func initTempLogger(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	Init(Options{
		ProjectName: name,
		LogDir:      dir,
		Level:       slog.LevelDebug,
		AddSource:   false,
	})
	logPath := path.Join(dir, name+".log")
	t.Logf("tmp log: %s", logPath)
	return logPath
}

// read entire log file
func readLogFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	return string(data)
}

func TestLogFunctions(t *testing.T) {
	logFile := initTempLogger(t, "app_test")

	// Call basic logging functions (they don't return values).
	Debug("debug message", slog.String("k", "v"))
	Info("info message", slog.Int("n", 1))
	Warn("warn message")
	Error("error message")

	Debugf("debugf %d", 1)
	Infof("infof %s", "x")
	Warnf("warnf %v", true)
	Errorf("errorf %0.2f", 3.14)

	// Notice KV
	PushNotice("a", 1)
	PushNotice("b", "two")
	Flush()

	// allow disk write
	time.Sleep(10 * time.Millisecond)

	content := readLogFile(t, logFile)

	if !strings.Contains(content, "debug message") {
		t.Errorf("missing debug message")
	}
	if !strings.Contains(content, "info message") {
		t.Errorf("missing info message")
	}
	if !strings.Contains(content, "warn message") {
		t.Errorf("missing warn message")
	}
	if !strings.Contains(content, "error message") {
		t.Errorf("missing error message")
	}
	if !strings.Contains(content, "NoticeKV") {
		t.Errorf("missing NoticeKV entry")
	}
}

// Test Fatal by running it in a separate process because it calls os.Exit(1)
func TestFatal(t *testing.T) {
	// compile-time path to this package's test binary
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalHelper")
	// set an env var so helper knows to run
	cmd.Env = append(os.Environ(), "GO_WANT_FATAL_TEST=1")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("expected process to exit non-zero")
	}
	content := readLogFile(t, string(output))
	if !strings.Contains(content, "fatal message") {
		t.Errorf("missing Fatal message")
	}
}

// helper that is executed in a subprocess
func TestFatalHelper(t *testing.T) {
	if os.Getenv("GO_WANT_FATAL_TEST") != "1" {
		return
	}
	fmt.Print(initTempLogger(t, "fatal_test"))
	// capture to stdout by writing to stderr as well
	Fatal("fatal message", "some error")
}
