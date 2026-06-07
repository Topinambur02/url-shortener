package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Host           string   `env:"HOST" env-default:"localhost" yaml:"host"`
		Port           int      `env:"PORT" env-default:"8080" yaml:"port"`
		StorageType    string   `env:"STORAGE_TYPE" env-default:"inmemory" yaml:"storage_type"`
		AllowedOrigins []string `env:"ALLOWED_ORIGINS" env-default:"*" yaml:"allowed_origins"`
	} `yaml:"app"`

	DB struct {
		Host     string `env:"DB_HOST" env-default:"localhost" yaml:"host"`
		Port     string `env:"DB_PORT" env-default:"5432" yaml:"port"`
		User     string `env:"DB_USER" env-default:"postgres" yaml:"user"`
		Password string `env:"DB_PASS" env-default:"postgres" yaml:"pass"`
		Name     string `env:"DB_NAME" env-default:"url_shortener_db" yaml:"name"`
	} `yaml:"db"`

	DSN string
}

func (c *Config) UpdateDSN() {
	c.DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Port,
	)
}

func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	fileInfo, err := os.Stat(filename)

	if err == nil && fileInfo.Size() > 0 {
		err = cleanenv.ReadConfig(filename, &cfg)

		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	err = cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to read environment: %w", err)
	}

	cfg.UpdateDSN()

	return &cfg, nil
}
