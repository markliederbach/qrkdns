package scheduler

import "github.com/go-co-op/gocron"

// Scheduler wraps the underlying gocron scheduler
type Scheduler interface {
	// Do specifies the jobFunc that should be called every time the Job runs
	Do(jobFun interface{}, params ...interface{}) (*gocron.Job, error)

	// StartBlocking starts all jobs and blocks the current thread
	StartBlocking()
}
