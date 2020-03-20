package resource

type Task struct {
	fn     func() error
	done   chan error
	errRun error
}

func (t *Task) Finished() bool {
	select {
	// FIXME
	case t.errRun = <-t.done:
		return true
	default:
		return false
	}
}

func (t *Task) Error() error {
	return t.errRun
}

func (t *Task) Wait() error {
	return <-t.done
}

func NewWorker() *Worker {
	return &Worker{
		queue: make(chan *Task, 16),
	}
}

type Worker struct {
	queue chan *Task
}

func (w *Worker) Work() {
	for task := range w.queue {
		err := task.fn()
		task.done <- err
	}
}

func (w *Worker) Schedule(fn func() error) *Task {
	task := &Task{
		fn:   fn,
		done: make(chan error, 1),
	}
	w.queue <- task
	return task
}
