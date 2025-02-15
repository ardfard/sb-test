package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	ID        uint      `db:"id"`
	AudioID   uint      `db:"audio_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SQLiteBroker struct {
	db *sqlx.DB
}

func NewSQLiteBroker(db *sqlx.DB) (*SQLiteBroker, error) {
	if err := createMessageTable(db); err != nil {
		return nil, fmt.Errorf("failed to create message table: %v", err)
	}
	return &SQLiteBroker{db: db}, nil
}

func createMessageTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS conversion_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		audio_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (audio_id) REFERENCES audios(id)
	)`
	_, err := db.Exec(query)
	return err
}

func (b *SQLiteBroker) EnqueueConversion(ctx context.Context, audioID uint) error {
	query := `
	INSERT INTO conversion_queue (
		audio_id, status, created_at, updated_at
	) VALUES (?, ?, ?, ?)`

	now := time.Now()
	_, err := b.db.ExecContext(ctx, query, audioID, "pending", now, now)
	if err != nil {
		return fmt.Errorf("failed to enqueue conversion: %v", err)
	}
	return nil
}

func (b *SQLiteBroker) DequeueConversion(ctx context.Context) (*Message, error) {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Get the oldest pending message and update its status atomically
	query := `
	SELECT id, audio_id, status, created_at, updated_at
	FROM conversion_queue
	WHERE status = 'pending'
	ORDER BY created_at ASC
	LIMIT 1`

	var msg Message
	err = tx.QueryRowContext(ctx, query).Scan(
		&msg.ID, &msg.AudioID, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %v", err)
	}

	// Update the message status to "processing"
	updateQuery := `
	UPDATE conversion_queue
	SET status = 'processing', updated_at = ?
	WHERE id = ?`

	now := time.Now()
	_, err = tx.ExecContext(ctx, updateQuery, now, msg.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update message status: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &msg, nil
}

func (b *SQLiteBroker) CompleteMessage(ctx context.Context, msgID uint) error {
	query := `
	UPDATE conversion_queue
	SET status = 'completed', updated_at = ?
	WHERE id = ?`

	_, err := b.db.ExecContext(ctx, query, time.Now(), msgID)
	if err != nil {
		return fmt.Errorf("failed to complete message: %v", err)
	}
	return nil
}
