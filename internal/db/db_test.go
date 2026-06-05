package db

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/topinambur02/url-shortener/internal/config"
)

func TestInitDBErrorCases(t *testing.T) {
	dsn := fmt.Sprintf("foo://bar:%s@qux", "baz")
	tests := []struct {
		name      string
		dsn       string
		wantError bool
		errSubstr string
	}{
		{
			name:      "empty DSN",
			dsn:       "",
			wantError: true,
			errSubstr: "failed to connect to database",
		},
		{
			name:      "invalid DSN format",
			dsn:       dsn,
			wantError: true,
			errSubstr: "failed to connect to database",
		},
		{
			name:      "non-existent host",
			dsn:       "host=non.existent.host port=5432 user=test dbname=test sslmode=disable",
			wantError: true,
			errSubstr: "failed to connect to database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := &config.Config{DSN: tt.dsn}
			db, err := InitDB(ctx, cfg)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !errors.Is(err, err) {
					if tt.errSubstr != "" && !contains(err.Error(), tt.errSubstr) {
						t.Errorf("error message %q does not contain %q", err.Error(), tt.errSubstr)
					}
				}

				if db != nil {
					t.Errorf("expected nil db on error, got %v", db)
				}

			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if db == nil {
					t.Error("expected non-nil db, got nil")
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}