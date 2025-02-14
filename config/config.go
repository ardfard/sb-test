package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds configuration values for the application.
type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	GCS           struct {
		Bucket string `mapstructure:"bucket"`
	} `mapstructure:"gcs"`
	SQLite struct {
		DBPath string `mapstructure:"db_path"`
	} `mapstructure:"sqlite"`
	Worker struct {
		NumWorkers int `mapstructure:"num_workers"`
	} `mapstructure:"worker"`
}

// LoadConfig reads configuration from config.yaml (or other supported formats) in the current directory.
func LoadConfig() (*Config, error) {
	// Set the name of the config file (without extension)
	viper.SetConfigName("config")
	// Set the config type (YAML in this example)
	viper.SetConfigType("yaml")
	// Look for the config file in the current directory.
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}
	return &cfg, nil
}
