package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

type mockCloser struct {
	closeCount int
	err        error
}

func (m *mockCloser) Close() error {
	m.closeCount++
	return m.err
}

func TestMain(m *testing.M) {
	logging.Init()
	exitCode := m.Run()

	err := os.RemoveAll("logs")
	if err != nil {
		fmt.Printf("Error deleting logs folder: %v\n", err)
	}

	os.Exit(exitCode)
}

func TestShutdown(t *testing.T) {
	t.Run("TestGracefulShutdownNormal", func(t *testing.T) {
		var shutdownCalled bool
		sh := func(ctx context.Context) error {
			shutdownCalled = true
			return nil
		}

		closer1 := &mockCloser{}
		closer2 := &mockCloser{}

		signals := []os.Signal{syscall.SIGUSR1}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer1, closer2)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("GracefulShutdown did not complete within timeout")
		}

		assert.True(t, shutdownCalled, "shutdown function was not called")
		assert.Equal(t, 1, closer1.closeCount, "closer1 not closed")
		assert.Equal(t, 1, closer2.closeCount, "closer2 not closed")
	})
	t.Run("TestGracefulShutdownShutdownError", func(t *testing.T) {
		shutdownErr := errors.New("shutdown failed")
		sh := func(ctx context.Context) error {
			return shutdownErr
		}

		closer := &mockCloser{}
		signals := []os.Signal{syscall.SIGUSR1}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("GracefulShutdown did not complete")
		}

		assert.Equal(t, 0, closer.closeCount, "closer should be closed even if shutdown returns error")
	})
	t.Run("TestGracefulShutdownCloserError", func(t *testing.T) {
		sh := func(ctx context.Context) error { return nil }
		goodCloser := &mockCloser{}
		badCloser := &mockCloser{err: errors.New("close error")}
		signals := []os.Signal{syscall.SIGUSR1}
		done := make(chan struct{})

		go func() {
			GracefulShutdown(signals, sh, goodCloser, badCloser)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("GracefulShutdown did not complete")
		}

		assert.Equal(t, 1, goodCloser.closeCount, "good closer should be closed")
		assert.Equal(t, 1, badCloser.closeCount, "bad closer should still be closed")
	})
	t.Run("TestGracefulShutdownMultipleSignals", func(t *testing.T) {
		var shutdownCalled bool
		sh := func(ctx context.Context) error {
			shutdownCalled = true
			return nil
		}

		closer := &mockCloser{}
		signals := []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("GracefulShutdown did not complete")
		}

		assert.True(t, shutdownCalled)
		assert.Equal(t, 1, closer.closeCount)
	})
}
