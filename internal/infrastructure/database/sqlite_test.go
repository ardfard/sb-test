package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name          string
		dbPath        string
		setup         func(string) error
		expectedError bool
	}{
		{
			name:          "Success - New Database",
			dbPath:        "test_new.db",
			expectedError: false,
		},
		{
			name:   "Success - Existing Database",
			dbPath: "test_existing.db",
			setup: func(path string) error {
				db, err := InitDB(path)
				if err != nil {
					return err
				}
				return db.Close()
			},
			expectedError: false,
		},
		{
			name:          "Error - Invalid Path",
			dbPath:        "/nonexistent/directory/test.db",
			expectedError: true,
		},
		{
			name:   "Error - Corrupted Database",
			dbPath: "test_corrupted.db",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("corrupted data"), 0644)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temp directory for test
			tempDir := t.TempDir()
			dbPath := filepath.Join(tempDir, tt.dbPath)

			// Run setup if provided
			if tt.setup != nil {
				err := tt.setup(dbPath)
				require.NoError(t, err)
			}

			// Test database initialization
			db, err := InitDB(dbPath)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, db)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, db)

			// Verify schema
			var tableExists bool
			err = db.Get(&tableExists, `
				SELECT EXISTS (
					SELECT 1 
					FROM sqlite_master 
					WHERE type='table' AND name='audios'
				)
			`)
			assert.NoError(t, err)
			assert.True(t, tableExists)

			// Verify columns
			rows, err := db.Query(`PRAGMA table_info(audios)`)
			assert.NoError(t, err)
			defer rows.Close()

			var columns []string
			for rows.Next() {
				var (
					cid       int
					name      string
					dtype     string
					notnull   int
					dfltValue interface{}
					pk        int
				)
				err := rows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk)
				assert.NoError(t, err)
				columns = append(columns, name)
			}

			// Verify expected columns exist
			expectedColumns := []string{
				"id",
				"original_name",
				"current_format",
				"storage_path",
				"status",
				"created_at",
				"updated_at",
				"error",
				"user_id",
				"phrase_id",
			}
			assert.Subset(t, columns, expectedColumns)

			// Clean up
			err = db.Close()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTable_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		db            *sqlx.DB
		expectedError string
	}{
		{
			name:          "Nil Database",
			db:            nil,
			expectedError: "database connection is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createTable(tt.db)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
