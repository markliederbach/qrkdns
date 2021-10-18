package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

// DefaultClient implements the scheduler client
type DefaultClient struct {
	CronSchedule string
	Client       Scheduler
}

// LoadOption allows for modifying the client after it's created
type LoadOption func(client *DefaultClient) error

// NewClient returns a new scheduler client
func NewClient(cronSchedule string, opts ...LoadOption) (DefaultClient, error) {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Cron(cronSchedule)
	coreJob := scheduler.Jobs()[0]
	if coreJob.Error() != nil {
		return DefaultClient{}, coreJob.Error()
	}

	client := DefaultClient{
		CronSchedule: cronSchedule,
		Client:       scheduler,
	}
	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return DefaultClient{}, err
		}
	}
	return client, nil
}

// GetScheduler returns a pointer to the underlying scheduler
func (c *DefaultClient) GetScheduler() Scheduler {
	return c.Client
}
