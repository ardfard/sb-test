package entity

import "time"

type AudioStatus string

const (
	AudioStatusPending    AudioStatus = "pending"
	AudioStatusConverting AudioStatus = "converting"
	AudioStatusCompleted  AudioStatus = "completed"
	AudioStatusFailed     AudioStatus = "failed"
)

type Audio struct {
	ID            uint        `db:"id"`
	OriginalName  string      `db:"original_name"`
	CurrentFormat string      `db:"current_format"`
	StoragePath   string      `db:"storage_path"`
	Status        AudioStatus `db:"status"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	Error         string      `db:"error"`
	UserID        uint        `db:"user_id"`
	PhraseID      uint        `db:"phrase_id"`
}
