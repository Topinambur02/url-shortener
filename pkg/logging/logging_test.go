package logging

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func setupTestLogger(t *testing.T) {
	t.Helper()

	oldE := e
	oldWd, err := os.Getwd()
	require.NoError(t, err, "failed to get current working directory")

	tmpDir, err := os.MkdirTemp("", "logging_test")
	require.NoError(t, err, "failed to create temp directory")

	err = os.Chdir(tmpDir)
	require.NoError(t, err, "failed to change working directory")

	Init()

	t.Cleanup(func() {
		e = oldE
		_ = os.Chdir(oldWd)
		_ = os.RemoveAll(tmpDir)
	})
}

func TestWriterHook_Fire(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	hook := &writerHook{
		Writer:    []io.Writer{buf1, buf2},
		LogLevels: []logrus.Level{logrus.InfoLevel},
	}

	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.InfoLevel,
		Message: "test message",
	}

	err := hook.Fire(entry)
	
	require.NoError(t, err)
	require.Contains(t, buf1.String(), "test message", "buf1 should contain the message")
	require.Contains(t, buf2.String(), "test message", "buf2 should contain the message")
}

func TestWriterHook_Levels(t *testing.T) {
	levels := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
	hook := &writerHook{
		LogLevels: levels,
	}
	
	require.Equal(t, levels, hook.Levels())
}

func TestLogging_Init(t *testing.T) {
	setupTestLogger(t)

	require.NotNil(t, e, "global logger entry should not be nil after Init")

	_, err := os.Stat("logs/all.log")
	require.NoError(t, err, "log file should be created")

	testMsg := "test message from Init test"
	e.Info(testMsg)

	content, err := os.ReadFile("logs/all.log")
	require.NoError(t, err)
	require.Contains(t, string(content), testMsg)
	require.Equal(t, logrus.TraceLevel, e.Logger.Level, "default log level should be TraceLevel")
}

func TestLogging_GetLogger(t *testing.T) {
	setupTestLogger(t)

	logger := GetLogger()
	require.NotNil(t, logger.Entry, "GetLogger should return a valid wrapper")
	
	testMsg := "GetLogger test"
	logger.Info(testMsg)

	content, err := os.ReadFile("logs/all.log")
	require.NoError(t, err)
	require.Contains(t, string(content), testMsg)
}

func TestLogging_GetLoggerWithField(t *testing.T) {
	setupTestLogger(t)

	logger := GetLogger()
	key := "mykey"
	value := "myvalue"

	loggerWithField := logger.GetLoggerWithField(key, value)

	require.Equal(t, value, loggerWithField.Data[key], "new logger should contain the custom field")

	_, ok := logger.Data[key]
	require.False(t, ok, "original logger should NOT be mutated")

	testMsg := "field test"
	loggerWithField.Info(testMsg)

	content, err := os.ReadFile("logs/all.log")
	require.NoError(t, err)
	require.Contains(t, string(content), key+"="+value, "log output should contain the formatted field")
	require.Contains(t, string(content), testMsg)
}

func TestLogging_InitDirAlreadyExists(t *testing.T) {
	setupTestLogger(t)

	require.NotPanics(t, Init, "Init should not panic if logs directory already exists")
}
