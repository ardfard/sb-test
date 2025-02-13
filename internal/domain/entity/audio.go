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
	ID             string
	OriginalName   string
	OriginalFormat string
	StoragePath    string
	Status         AudioStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Error          string
}
