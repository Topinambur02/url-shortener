package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("TestLoadConfigFromFile", func(t *testing.T) {
		content := `
app:
  host: "testhost"
  port: 9090

db:
  host: "dbhost"
  port: "5433"
  user: "testuser"
  pass: "testpass"
  name: "testdb"
`
		tmpFile, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)
		_ = tmpFile.Close()

		cfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)

		require.Equal(t, "testhost", cfg.App.Host)
		require.Equal(t, 9090, cfg.App.Port)
		require.Equal(t, "dbhost", cfg.DB.Host)
		require.Equal(t, "5433", cfg.DB.Port)
		require.Equal(t, "testuser", cfg.DB.User)
		require.Equal(t, "testpass", cfg.DB.Password)
		require.Equal(t, "testdb", cfg.DB.Name)

		expectedDSN := "host=dbhost user=testuser password=testpass dbname=testdb port=5433 sslmode=disable"
		require.Equal(t, expectedDSN, cfg.DSN)
	})
	t.Run("TestLoadConfigFromEnvWhenFileMissing", func(t *testing.T) {
		_ = os.Setenv("HOST", "envhost")
		_ = os.Setenv("PORT", "9091")
		_ = os.Setenv("DB_HOST", "envdbhost")
		_ = os.Setenv("DB_PORT", "5434")
		_ = os.Setenv("DB_USER", "envuser")
		_ = os.Setenv("DB_PASS", "envpass")
		_ = os.Setenv("DB_NAME", "envdb")

		defer func() {
			os.Unsetenv("HOST")
			os.Unsetenv("PORT")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASS")
			os.Unsetenv("DB_NAME")
		}()

		cfg, err := LoadConfig("nonexistent.yaml")
		require.NoError(t, err)

		require.Equal(t, "envhost", cfg.App.Host)
		require.Equal(t, 9091, cfg.App.Port)
		require.Equal(t, "envdbhost", cfg.DB.Host)
		require.Equal(t, "5434", cfg.DB.Port)
		require.Equal(t, "envuser", cfg.DB.User)
		require.Equal(t, "envpass", cfg.DB.Password)
		require.Equal(t, "envdb", cfg.DB.Name)

		expectedDSN := "host=envdbhost user=envuser password=envpass dbname=envdb port=5434 sslmode=disable"
		require.Equal(t, expectedDSN, cfg.DSN)
	})
	t.Run("TestLoadConfigDefaultsWhenNoFileAndNoEnv", func(t *testing.T) {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")

		cfg, err := LoadConfig("nonexistent.yaml")
		require.NoError(t, err)

		require.Equal(t, "localhost", cfg.App.Host)
		require.Equal(t, 8080, cfg.App.Port)
		require.Equal(t, "localhost", cfg.DB.Host)
		require.Equal(t, "5432", cfg.DB.Port)
		require.Equal(t, "postgres", cfg.DB.User)
		require.Equal(t, "postgres", cfg.DB.Password)
		require.Equal(t, "url_shortener_db", cfg.DB.Name)

		expectedDSN := "host=localhost user=postgres password=postgres dbname=url_shortener_db port=5432 sslmode=disable"
		require.Equal(t, expectedDSN, cfg.DSN)
	})
	t.Run("TestLoadConfigEnvOverridesFile", func(t *testing.T) {
		content := `
app:
  host: "filehost"
  port: 8081

db:
  host: "filedbhost"
  port: "5435"
  user: "fileuser"
  pass: "filepass"
  name: "filedb"
`
		tmpFile, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)
		_ = tmpFile.Close()

		_ = os.Setenv("APP_HOST", "envhost")
		_ = os.Setenv("HOST", "envhost")
		_ = os.Setenv("PORT", "9092")
		_ = os.Setenv("DB_HOST", "envdbhost")
		_ = os.Setenv("DB_PASS", "envpass")

		defer func() {
			os.Unsetenv("HOST")
			os.Unsetenv("PORT")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PASS")
		}()

		cfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)

		require.Equal(t, "envhost", cfg.App.Host)
		require.Equal(t, 9092, cfg.App.Port)
		require.Equal(t, "envdbhost", cfg.DB.Host)
		require.Equal(t, "5435", cfg.DB.Port)
		require.Equal(t, "fileuser", cfg.DB.User)
		require.Equal(t, "envpass", cfg.DB.Password)
		require.Equal(t, "filedb", cfg.DB.Name)

		expectedDSN := "host=envdbhost user=fileuser password=envpass dbname=filedb port=5435 sslmode=disable"
		require.Equal(t, expectedDSN, cfg.DSN)
	})
	t.Run("TestLoadConfigInvalidFileReturnsError", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("this is not valid yaml: [")
		require.NoError(t, err)
		_ = tmpFile.Close()

		cfg, err := LoadConfig(tmpFile.Name())
		require.Error(t, err)
		require.Nil(t, cfg)
	})
	t.Run("TestLoadConfigEmptyFileUsesDefaults", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("")
		require.NoError(t, err)
		_ = tmpFile.Close()

		cfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)

		require.Equal(t, "localhost", cfg.App.Host)
		require.Equal(t, 8080, cfg.App.Port)
		require.Equal(t, "localhost", cfg.DB.Host)
		require.Equal(t, "5432", cfg.DB.Port)
		require.Equal(t, "postgres", cfg.DB.User)
		require.Equal(t, "postgres", cfg.DB.Password)
		require.Equal(t, "url_shortener_db", cfg.DB.Name)
	})
	t.Run("TestUpdateDSN", func(t *testing.T) {
		cfg := &Config{}
		cfg.DB.Host = "testhost"
		cfg.DB.User = "testuser"
		cfg.DB.Password = "testpass"
		cfg.DB.Name = "testdb"
		cfg.DB.Port = "1234"

		cfg.UpdateDSN()

		expected := "host=testhost user=testuser password=testpass dbname=testdb port=1234 sslmode=disable"
		require.Equal(t, expected, cfg.DSN)
	})
}
