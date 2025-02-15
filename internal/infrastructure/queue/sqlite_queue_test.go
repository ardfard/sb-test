package queue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}

func TestNewSQLiteQueue(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("successful initialization", func(t *testing.T) {
		queue, err := NewSQLiteQueue(db, "audio_conversion")
		require.NoError(t, err)
		assert.NotNil(t, queue)
	})

	t.Run("multiple initializations should not error", func(t *testing.T) {
		// First initialization
		_, err := NewSQLiteQueue(db, "audio_conversion")
		require.NoError(t, err)

		// Second initialization should not error
		queue2, err := NewSQLiteQueue(db, "audio_conversion")
		require.NoError(t, err)
		assert.NotNil(t, queue2)
	})
}

func TestSQLiteQueueOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("enqueue and dequeue single task", func(t *testing.T) {
		queueName := fmt.Sprintf("test_queue_%s_single", t.Name())
		queue, err := NewSQLiteQueue(db, queueName)
		require.NoError(t, err)

		// Enqueue a task
		payload := uint(123)
		err = queue.Enqueue(ctx, payload)
		require.NoError(t, err)

		// Dequeue the task
		task, err := queue.Dequeue(ctx)
		require.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, payload, task.Payload)
	})

	t.Run("complete task", func(t *testing.T) {
		queueName := fmt.Sprintf("test_queue_%s_complete", t.Name())
		queue, err := NewSQLiteQueue(db, queueName)
		require.NoError(t, err)

		// Enqueue a task
		payload := uint(456)
		err = queue.Enqueue(ctx, payload)
		require.NoError(t, err)

		// Dequeue the task
		task, err := queue.Dequeue(ctx)
		require.NoError(t, err)

		// Complete the task
		err = queue.Complete(ctx, task.ID)
		require.NoError(t, err)

		// Try to dequeue again - should timeout as queue is empty
		timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		_, err = queue.Dequeue(timeoutCtx)
		assert.Error(t, err) // Should error due to timeout
	})

	t.Run("fail task", func(t *testing.T) {
		queueName := fmt.Sprintf("test_queue_%s_fail", t.Name())
		// Enqueue a task
		queue, err := NewSQLiteQueue(db, queueName)
		require.NoError(t, err)

		payload := uint(789)
		err = queue.Enqueue(ctx, payload)
		require.NoError(t, err)

		// Dequeue the task
		task, err := queue.Dequeue(ctx)
		require.NoError(t, err)

		// Fail the task
		err = queue.Fail(ctx, task.ID, "test error message")
		require.NoError(t, err)

		// Try to dequeue again - should timeout as queue is empty
		timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		_, err = queue.Dequeue(timeoutCtx)
		assert.Error(t, err) // Should error due to timeout
	})

	t.Run("dequeue from empty queue", func(t *testing.T) {
		queueName := fmt.Sprintf("test_queue_%s_empty", t.Name())
		queue, err := NewSQLiteQueue(db, queueName)
		// Try to dequeue from empty queue with short timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		_, err = queue.Dequeue(timeoutCtx)
		assert.Error(t, err) // Should error due to timeout
	})

	t.Run("multiple enqueue and dequeue", func(t *testing.T) {
		queueName := fmt.Sprintf("test_queue_%s_multiple", t.Name())
		queue, err := NewSQLiteQueue(db, queueName)
		payloads := []uint{1, 2, 3, 4, 5}

		// Enqueue multiple tasks
		for _, p := range payloads {
			err = queue.Enqueue(ctx, p)
			require.NoError(t, err)
		}

		// Dequeue and verify all tasks
		for _, expectedPayload := range payloads {
			task, err := queue.Dequeue(ctx)
			require.NoError(t, err)
			assert.Equal(t, expectedPayload, task.Payload)

			// Complete the task
			err = queue.Complete(ctx, task.ID)
			require.NoError(t, err)
		}
	})
}
