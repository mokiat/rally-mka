package resource

import "log"

type Task func() error

func NewWorker(queueSize int) *Worker {
	return &Worker{
		queue: make(chan Task, queueSize),
	}
}

type Worker struct {
	queue chan Task
}

func (w *Worker) Schedule(task Task) {
	w.queue <- task
}

func (w *Worker) Work() {
	for task := range w.queue {
		if err := task(); err != nil {
			log.Printf("failed to execute task: %v", err)
			return
		}
	}
}
