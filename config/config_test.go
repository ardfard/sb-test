package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server_address: ":8080"
storage:
  type: "s3"
  s3:
    bucket: "test-bucket"
    region: "us-east-1"
    access_key_id: "test-key"
    secret_access_key: "test-secret"
  local:
    directory: "/tmp/storage"
sqlite:
  db_path: "/tmp/app.db"
worker:
  num_workers: 4
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		configPath  string
		wantErr     bool
		validateCfg func(*testing.T, *Config)
	}{
		{
			name:       "valid config",
			configPath: configPath,
			wantErr:    false,
			validateCfg: func(t *testing.T, cfg *Config) {
				assert.Equal(t, ":8080", cfg.ServerAddress)
				assert.Equal(t, "s3", cfg.Storage.Type)
				assert.Equal(t, "test-bucket", cfg.Storage.S3.Bucket)
				assert.Equal(t, "us-east-1", cfg.Storage.S3.Region)
				assert.Equal(t, "test-key", cfg.Storage.S3.AccessKeyID)
				assert.Equal(t, "test-secret", cfg.Storage.S3.SecretAccessKey)
				assert.Equal(t, "/tmp/storage", cfg.Storage.Local.Directory)
				assert.Equal(t, "/tmp/app.db", cfg.SQLite.DBPath)
			},
		},
		{
			name:       "non-existent config file",
			configPath: "non-existent.yaml",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.configPath)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
			if tt.validateCfg != nil {
				tt.validateCfg(t, cfg)
			}
		})
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")

	invalidContent := `
server_address: ":8080"
storage:
  type: "s3"
  s3:
    bucket: [invalid yaml
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
