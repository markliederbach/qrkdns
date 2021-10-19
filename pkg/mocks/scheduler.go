package mocks

import (
	"github.com/go-co-op/gocron"
	"github.com/markliederbach/go-envy"
)

var (
	// DefaultSchedulerDoResponse is the default response for this function
	DefaultSchedulerDoResponse *gocron.Job = &gocron.Job{}
)

// MockSchedulerClient mocks the internal scheduler
type MockSchedulerClient struct{}

func init() {
	sdkFunctions := []string{
		"Do",
		"StartBlocking",
	}
	for _, functionName := range sdkFunctions {
		envy.ObjectChannels[functionName] = make(chan interface{}, 100)
		envy.ErrorChannels[functionName] = make(chan error, 100)
		envy.DefaultObjects[functionName] = struct{}{}
		envy.DefaultErrors[functionName] = nil
	}
}

// Do implements corresponding client function
func (c *MockSchedulerClient) Do(jobFun interface{}, params ...interface{}) (*gocron.Job, error) {
	functionName := "Do"
	obj := envy.GetObject(functionName)
	err := envy.GetError(functionName)
	switch obj := obj.(type) {
	case *gocron.Job:
		return obj, err
	default:
		return DefaultSchedulerDoResponse, err
	}

}

// StartBlocking implements corresponding client function
func (c *MockSchedulerClient) StartBlocking() {}
