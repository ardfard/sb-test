package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/infrastructure/database"
)

func TestSQLiteAudioRepository(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*AudioRepository) (*entity.Audio, error)
		check   func(*testing.T, *AudioRepository, *entity.Audio)
		wantErr bool
	}{
		{
			name: "Store and retrieve audio",
			setup: func(repo *AudioRepository) (*entity.Audio, error) {
				audio := &entity.Audio{
					ID:            1,
					OriginalName:  "test.m4a",
					CurrentFormat: ".m4a",
					StoragePath:   "original/test.m4a",
					Status:        entity.AudioStatusPending,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
					UserID:        1,
					PhraseID:      1,
				}
				err := repo.Store(context.Background(), audio)
				return audio, err
			},
			check: func(t *testing.T, repo *AudioRepository, audio *entity.Audio) {
				stored, err := repo.GetByID(context.Background(), audio.ID)
				if err != nil {
					t.Fatalf("GetByID failed: %v", err)
				}
				if stored.ID != audio.ID || stored.OriginalName != audio.OriginalName {
					t.Errorf("Stored audio mismatch:\ngot  %+v\nwant %+v", stored, audio)
				}
			},
			wantErr: false,
		},
		{
			name: "Update audio status",
			setup: func(repo *AudioRepository) (*entity.Audio, error) {
				audio := &entity.Audio{
					ID:            2,
					OriginalName:  "test2.m4a",
					CurrentFormat: ".m4a",
					StoragePath:   "original/test2.m4a",
					Status:        entity.AudioStatusPending,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
					UserID:        1,
					PhraseID:      1,
				}
				if err := repo.Store(context.Background(), audio); err != nil {
					return nil, err
				}
				audio.Status = entity.AudioStatusCompleted
				audio.StoragePath = "converted/test2.wav"
				err := repo.Update(context.Background(), audio)
				return audio, err
			},
			check: func(t *testing.T, repo *AudioRepository, audio *entity.Audio) {
				updated, err := repo.GetByID(context.Background(), audio.ID)
				if err != nil {
					t.Fatalf("GetByID after update failed: %v", err)
				}
				if updated.Status != entity.AudioStatusCompleted || updated.StoragePath != "converted/test2.wav" {
					t.Errorf("Updated audio mismatch:\ngot  %+v\nwant status=%s, storagePath=%s",
						updated, entity.AudioStatusCompleted, "converted/test2.wav")
				}
			},
			wantErr: false,
		},
		{
			name: "Get non-existent audio",
			setup: func(repo *AudioRepository) (*entity.Audio, error) {
				return &entity.Audio{ID: 9999}, nil
			},
			check: func(t *testing.T, repo *AudioRepository, audio *entity.Audio) {
				_, err := repo.GetByID(context.Background(), audio.ID)
				if err == nil {
					t.Error("expected error when fetching nonexistent audio, got nil")
				}
			},
			wantErr: false,
		},
		{
			name: "Update non-existent audio",
			setup: func(repo *AudioRepository) (*entity.Audio, error) {
				audio := &entity.Audio{
					ID:            9999,
					OriginalName:  "nonexistent.m4a",
					CurrentFormat: ".m4a",
					Status:        entity.AudioStatusCompleted,
				}
				err := repo.Update(context.Background(), audio)
				return audio, err
			},
			check: func(t *testing.T, repo *AudioRepository, audio *entity.Audio) {
				// No additional checks needed as we expect the setup to fail
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using in-memory SQLite DB for each test
			db, err := database.InitDB(":memory:")
			if err != nil {
				t.Fatalf("failed to create repository: %v", err)
			}
			repo, err := NewAudioRepository(db)
			if err != nil {
				t.Fatalf("failed to create repository: %v", err)
			}
			defer repo.Close()

			audio, err := tt.setup(repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("setup error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				tt.check(t, repo, audio)
			}
		})
	}
}
