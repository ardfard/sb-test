package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ardfard/sb-test/internal/domain/queue"
	"github.com/ardfard/sb-test/internal/usecase"
)

// ConversionWorker is a worker that converts audio files.
type ConversionWorker struct {
	queue    queue.TaskQueue
	useCase  *usecase.ConvertAudioUseCase
	stopChan chan struct{}
}

// NewConversionWorker creates a new ConversionWorker.
func NewConversionWorker(
	queue queue.TaskQueue,
	convertUseCase *usecase.ConvertAudioUseCase,
) *ConversionWorker {
	return &ConversionWorker{
		queue:    queue,
		useCase:  convertUseCase,
		stopChan: make(chan struct{}),
	}
}

// Start starts the worker in a new goroutine.
func (w *ConversionWorker) Start() {
	go w.run()
}

// Stop stops the worker.
func (w *ConversionWorker) Stop() {
	close(w.stopChan)
}

// run is the main loop for the worker.
// It continuously processes messages from the queue.
func (w *ConversionWorker) run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			if err := w.processNextMessage(); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}

// processNextMessage processes the next message from the queue.
func (w *ConversionWorker) processNextMessage() (err error) {
	ctx := context.Background()

	task, err := w.queue.Dequeue(ctx)
	defer func() {
		if err != nil && task != nil {
			if err := w.queue.Fail(ctx, task.ID, err.Error()); err != nil {
				log.Printf("failed to fail task: %v", err)
			}
		}
	}()

	if err != nil {
		err = fmt.Errorf("failed to dequeue message: %v", err)
		return
	}

	if err = w.useCase.Convert(ctx, task.Payload); err != nil {
		err = fmt.Errorf("failed to convert audio: %v", err)
		return
	}

	// Mark message as completed
	if err := w.queue.Complete(ctx, task.ID); err != nil {
		return fmt.Errorf("failed to complete message: %v", err)
	}

	return nil
}
