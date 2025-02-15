package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// StorageConfig holds settings for the storage backend.
type StorageConfig struct {
	Type string `mapstructure:"type"` // Allowed values: "gcs", "s3", "local"
	S3   struct {
		Bucket          string `mapstructure:"bucket"`
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"s3"`
	Local struct {
		Directory string `mapstructure:"directory"`
	} `mapstructure:"local"`
}

// Config holds configuration values for the application.
type Config struct {
	ServerAddress string        `mapstructure:"server_address"`
	Storage       StorageConfig `mapstructure:"storage"`
	SQLite        struct {
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
