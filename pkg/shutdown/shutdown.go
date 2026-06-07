package shutdown

import (
	"context"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/topinambur02/url-shortener/pkg/logging"
)

func GracefulShutdown(signals []os.Signal, shutdown func(context.Context) error, closers ...io.Closer) error {
	logger := logging.GetLogger()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	<-sigc
	logger.Info("Caught signal. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		return err
	}

	for _, c := range closers {
		if err := c.Close(); err != nil {
			return err
		}
	}

	return nil
}
