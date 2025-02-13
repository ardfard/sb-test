package repository

import (
	"context"
	"testing"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"
)

func TestSQLiteAudioRepository(t *testing.T) {
	ctx := context.Background()

	// Using in-memory SQLite DB for testing.
	repo, err := NewSQLiteAudioRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	now := time.Now().UTC()
	audio := &entity.Audio{
		ID:             "test-audio-id",
		OriginalName:   "test.wav",
		OriginalFormat: ".wav",
		Status:         entity.AudioStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
		Error:          "",
	}

	// Test storing an audio record.
	if err := repo.Store(ctx, audio); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Test retrieving the stored audio.
	stored, err := repo.GetByID(ctx, audio.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if stored.ID != audio.ID || stored.OriginalName != audio.OriginalName || stored.Status != audio.Status {
		t.Errorf("Stored audio mismatch:\ngot  %+v\nwant %+v", stored, audio)
	}

	// Test updating the audio record.
	audio.Status = entity.AudioStatusCompleted
	audio.StoragePath = "converted/test.wav"
	if err := repo.Update(ctx, audio); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.GetByID(ctx, audio.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if updated.Status != entity.AudioStatusCompleted || updated.StoragePath != "converted/test.wav" {
		t.Errorf("Updated audio mismatch:\ngot  %+v\nwant status=%s, storagePath=%s", updated, entity.AudioStatusCompleted, "converted/test.wav")
	}
}

func TestGetNonExistentAudio(t *testing.T) {
	ctx := context.Background()
	repo, err := NewSQLiteAudioRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	_, err = repo.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error when fetching nonexistent audio, got nil")
	}
}
