// Package queue provides queue implementations for task processing.
// This package implements the queue.TaskQueue interface using SQLite as the underlying storage.
// It uses the goqite library for reliable queue operations and maintains a separate table
// for tracking failed tasks.
package queue

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ardfard/sb-test/internal/domain/queue"
	"github.com/jmoiron/sqlx"
	"github.com/maragudk/goqite"
)

// SQLiteQueue implements queue.TaskQueue using goqite
type SQLiteQueue struct {
	queue *goqite.Queue
	db    *sqlx.DB // Used for failed tasks table operations
}

// NewSQLiteQueue creates a new SQLiteQueue
func NewSQLiteQueue(db *sqlx.DB, queueName string) (*SQLiteQueue, error) {
	// Convert sqlx.DB to sql.DB since goqite expects sql.DB
	sqlDB := db.DB

	// Try to initialize goqite schema, ignore "table already exists" error
	err := goqite.Setup(context.Background(), db.DB)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, fmt.Errorf("failed to initialize queue schema: %v", err)
	}

	// Create a new queue named "audio_conversion" with default settings
	q := goqite.New(goqite.NewOpts{
		DB:   sqlDB,
		Name: queueName,
	})

	return &SQLiteQueue{
		queue: q,
		db:    db,
	}, nil
}

func (q *SQLiteQueue) Enqueue(ctx context.Context, payload uint) error {
	err := q.queue.Send(ctx, goqite.Message{
		Body: []byte(fmt.Sprintf("%d", payload)),
	})
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %v", err)
	}
	return nil
}

// Dequeue gets a task from the queue
func (q *SQLiteQueue) Dequeue(ctx context.Context) (*queue.Task, error) {
	msg, err := q.queue.ReceiveAndWait(ctx, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue message: %v", err)
	}

	// Convert the message body back to uint
	var payload uint
	_, err = fmt.Sscanf(string(msg.Body), "%d", &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message payload: %v", err)
	}

	return &queue.Task{
		ID:      string(msg.ID),
		Payload: payload,
	}, nil
}

// Complete marks a task as completed
func (q *SQLiteQueue) Complete(ctx context.Context, taskID string) error {
	err := q.queue.Delete(ctx, goqite.ID(taskID))
	if err != nil {
		return fmt.Errorf("failed to complete message: %v", err)
	}
	return nil
}

// Fail marks a task as failed, stores it in the failed_tasks table, and removes it from the queue
func (q *SQLiteQueue) Fail(ctx context.Context, taskID string, errMsg string) error {
	// Delete and do nothing for now
	err := q.queue.Delete(ctx, goqite.ID(taskID))
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
