package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/topinambur02/url-shortener/pkg/logging"
	"github.com/topinambur02/url-shortener/pkg/shutdown/mocks"
)

func TestMain(m *testing.M) {
	logging.Init()
	exitCode := m.Run()

	if err := os.RemoveAll("logs"); err != nil {
		fmt.Printf("Error deleting logs folder: %v\n", err)
	}

	os.Exit(exitCode)
}

func TestGracefulShutdown(t *testing.T) {
	triggerAndAwait := func(t *testing.T, done chan struct{}, sig syscall.Signal) {
		t.Helper()
		
		time.Sleep(50 * time.Millisecond)
		
		err := syscall.Kill(syscall.Getpid(), sig)
		require.NoError(t, err, "failed to send signal")

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			require.FailNow(t, "GracefulShutdown did not complete within timeout")
		}
	}

	t.Run("normal execution", func(t *testing.T) {
		var shutdownCalled bool
		sh := func(ctx context.Context) error {
			shutdownCalled = true
			return nil
		}

		closer1 := &mocks.MockCloser{}
		closer2 := &mocks.MockCloser{}
		signals := []os.Signal{syscall.SIGUSR1}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer1, closer2)
			close(done)
		}()

		triggerAndAwait(t, done, syscall.SIGUSR1)

		require.True(t, shutdownCalled, "shutdown function was not called")
		require.Equal(t, 1, closer1.CloseCount, "closer1 not closed")
		require.Equal(t, 1, closer2.CloseCount, "closer2 not closed")
	})

	t.Run("shutdown error", func(t *testing.T) {
		shutdownErr := errors.New("shutdown failed")
		sh := func(ctx context.Context) error {
			return shutdownErr
		}

		closer := &mocks.MockCloser{}
		signals := []os.Signal{syscall.SIGUSR1}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer)
			close(done)
		}()

		triggerAndAwait(t, done, syscall.SIGUSR1)

		require.Equal(t, 0, closer.CloseCount, "closer should NOT be closed if shutdown returns error")
	})

	t.Run("closer error", func(t *testing.T) {
		sh := func(ctx context.Context) error { return nil }
		goodCloser := &mocks.MockCloser{}
		badCloser := &mocks.MockCloser{Err: errors.New("close error")}
		signals := []os.Signal{syscall.SIGUSR1}
		
		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, goodCloser, badCloser)
			close(done)
		}()

		triggerAndAwait(t, done, syscall.SIGUSR1)

		require.Equal(t, 1, goodCloser.CloseCount, "good closer should be closed")
		require.Equal(t, 1, badCloser.CloseCount, "bad closer should still be closed despite error")
	})

	t.Run("multiple signals", func(t *testing.T) {
		var shutdownCalled bool
		sh := func(ctx context.Context) error {
			shutdownCalled = true
			return nil
		}

		closer := &mocks.MockCloser{}
		signals := []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}

		done := make(chan struct{})
		go func() {
			GracefulShutdown(signals, sh, closer)
			close(done)
		}()

		triggerAndAwait(t, done, syscall.SIGUSR2)

		require.True(t, shutdownCalled, "shutdown function was not called")
		require.Equal(t, 1, closer.CloseCount, "closer not closed")
	})
}
