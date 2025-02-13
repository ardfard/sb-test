package worker

import (
	"sync"
)

type Job func()

type Worker struct {
	jobQueue chan Job
	wg       sync.WaitGroup
}

func NewWorker(numWorkers int) *Worker {
	w := &Worker{
		jobQueue: make(chan Job, 100),
	}

	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.work()
	}

	return w
}

func (w *Worker) work() {
	defer w.wg.Done()
	for job := range w.jobQueue {
		job()
	}
}

func (w *Worker) EnqueueJob(job Job) {
	w.jobQueue <- job
}

func (w *Worker) Shutdown() {
	close(w.jobQueue)
	w.wg.Wait()
}
