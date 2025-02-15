package queue

import (
	"context"
	"time"
)

// Task represents a unit of work to be processed
type Task struct {
	ID        string    // Unique identifier for the task
	Type      string    // Type of task (e.g., "audio_conversion")
	Payload   uint      // Task-specific data (in this case, audio ID)
	Status    string    // Current status of the task
	CreatedAt time.Time // When the task was created
	UpdatedAt time.Time // When the task was last updated
}

// TaskQueue defines the interface for a task queue implementation
type TaskQueue interface {
	// Enqueue adds a new task to the queue
	Enqueue(ctx context.Context, payload uint) error

	// Dequeue retrieves and claims the next available task
	Dequeue(ctx context.Context) (*Task, error)

	// Complete marks a task as completed
	Complete(ctx context.Context, taskID string) error

	// Fail marks a task as failed with an error message
	Fail(ctx context.Context, taskID string, errMsg string) error
}
