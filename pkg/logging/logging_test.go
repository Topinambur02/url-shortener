package logging

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupTestLogger(t *testing.T) func() {
	oldE := e
	oldWd, err := os.Getwd()
	assert.NoError(t, err)

	tmpDir, err := os.MkdirTemp("", "logging_test")
	assert.NoError(t, err)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	Init()

	cleanup := func() {
		e = oldE
		_ = os.Chdir(oldWd)
		_ = os.RemoveAll(tmpDir)
	}
	return cleanup
}

func TestLogging(t *testing.T) {
	t.Run("WriterHookFire", func(t *testing.T) {
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
		assert.NoError(t, err)
		assert.Contains(t, buf1.String(), "test message")
		assert.Contains(t, buf2.String(), "test message")
	})
	t.Run("WriterHookLevels", func(t *testing.T) {
		levels := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
		hook := &writerHook{
			LogLevels: levels,
		}
		assert.Equal(t, levels, hook.Levels())
	})
	t.Run("Init", func(t *testing.T) {
		cleanup := setupTestLogger(t)
		defer cleanup()

		assert.NotNil(t, e)

		_, err := os.Stat("logs/all.log")
		assert.NoError(t, err)

		testMsg := "test message from Init test"
		e.Info(testMsg)

		content, err := os.ReadFile("logs/all.log")
		assert.NoError(t, err)
		assert.Contains(t, string(content), testMsg)
		assert.Equal(t, logrus.TraceLevel, e.Logger.Level)
	})
	t.Run("GetLogger", func(t *testing.T) {
		cleanup := setupTestLogger(t)
		defer cleanup()

		logger := GetLogger()
		assert.NotNil(t, logger.Entry)
		testMsg := "GetLogger test"
		logger.Info(testMsg)

		content, err := os.ReadFile("logs/all.log")
		assert.NoError(t, err)
		assert.Contains(t, string(content), testMsg)
	})
	t.Run("GetLoggerWithField", func(t *testing.T) {
		cleanup := setupTestLogger(t)
		defer cleanup()

		logger := GetLogger()
		key := "mykey"
		value := "myvalue"

		loggerWithField := logger.GetLoggerWithField(key, value)

		assert.Equal(t, value, loggerWithField.Data[key])

		_, ok := logger.Data[key]
		assert.False(t, ok)

		testMsg := "field test"
		loggerWithField.Info(testMsg)

		content, err := os.ReadFile("logs/all.log")
		assert.NoError(t, err)
		assert.Contains(t, string(content), key+"="+value)
		assert.Contains(t, string(content), testMsg)
	})
	t.Run("InitDirAlreadyExists", func(t *testing.T) {
		cleanup := setupTestLogger(t)
		defer cleanup()

		assert.NotPanics(t, Init)
	})
}
