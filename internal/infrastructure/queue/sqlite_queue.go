package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/ardfard/sb-test/internal/domain/queue"
	"github.com/jmoiron/sqlx"
)

type SQLiteQueue struct {
	db *sqlx.DB
}

func NewSQLiteQueue(db *sqlx.DB) (*SQLiteQueue, error) {
	if err := createTaskTable(db); err != nil {
		return nil, fmt.Errorf("failed to create task table: %v", err)
	}
	return &SQLiteQueue{db: db}, nil
}

func createTaskTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		payload INTEGER NOT NULL,
		status TEXT NOT NULL,
		error TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func (q *SQLiteQueue) Enqueue(ctx context.Context, taskType string, payload uint) error {
	query := `
	INSERT INTO tasks (
		type, payload, status, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	_, err := q.db.ExecContext(ctx, query, taskType, payload, "pending", now, now)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}
	return nil
}

func (q *SQLiteQueue) Dequeue(ctx context.Context, taskType string) (*queue.Task, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := `
	SELECT id, type, payload, status, created_at, updated_at
	FROM tasks
	WHERE type = ? AND status = 'pending'
	ORDER BY created_at ASC
	LIMIT 1`

	var task queue.Task
	err = tx.QueryRowContext(ctx, query, taskType).Scan(
		&task.ID,
		&task.Type,
		&task.Payload,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	updateQuery := `
	UPDATE tasks
	SET status = 'processing', updated_at = ?
	WHERE id = ?`

	now := time.Now()
	_, err = tx.ExecContext(ctx, updateQuery, now, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &task, nil
}

func (q *SQLiteQueue) Complete(ctx context.Context, taskID uint) error {
	query := `
	UPDATE tasks
	SET status = 'completed', updated_at = ?
	WHERE id = ?`

	_, err := q.db.ExecContext(ctx, query, time.Now(), taskID)
	if err != nil {
		return fmt.Errorf("failed to complete task: %v", err)
	}
	return nil
}

func (q *SQLiteQueue) Fail(ctx context.Context, taskID uint, errMsg string) error {
	query := `
	UPDATE tasks
	SET status = 'failed', error = ?, updated_at = ?
	WHERE id = ?`

	_, err := q.db.ExecContext(ctx, query, errMsg, time.Now(), taskID)
	if err != nil {
		return fmt.Errorf("failed to mark task as failed: %v", err)
	}
	return nil
}
