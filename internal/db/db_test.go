package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/topinambur02/url-shortener/internal/config"
)

func TestInitDBErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dsn       string
		errSubstr string
	}{
		{
			name:      "empty DSN",
			dsn:       "",
			errSubstr: "failed to connect to database",
		},
		{
			name:      "invalid DSN format",
			dsn:       "foo://bar:baz@qux",
			errSubstr: "unable to create connection pool",
		},
		{
			name:      "non-existent host",
			dsn:       "host=non.existent.host port=5432 user=test dbname=test sslmode=disable",
			errSubstr: "failed to connect to database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := &config.Config{DSN: tt.dsn}
			db, err := InitDB(ctx, cfg)

			require.Error(t, err, "An error was expected but was received nil")
			require.Contains(t, err.Error(), tt.errSubstr, "The error message %q does not contain %q", err.Error(), tt.errSubstr)
			require.Nil(t, db, "in case of error DB should be nil")
		})
	}
}
