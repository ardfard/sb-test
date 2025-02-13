package worker

import (
	"sync"
	"testing"
)

func TestWorkerEnqueue(t *testing.T) {
	w := NewWorker(3)
	var mu sync.Mutex
	counter := 0

	totalJobs := 10
	for i := 0; i < totalJobs; i++ {
		w.EnqueueJob(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	// Shutdown the worker so all jobs are processed.
	w.Shutdown()

	if counter != totalJobs {
		t.Errorf("expected counter to be %d, got %d", totalJobs, counter)
	}
}
