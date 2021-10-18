package scheduler_test

import (
	"fmt"
	"testing"

	"github.com/markliederbach/qrkdns/pkg/clients/scheduler"
	"github.com/markliederbach/qrkdns/pkg/mocks"
	. "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func newMockSchedulerClient() (scheduler.DefaultClient, error) {
	client, err := scheduler.NewClient("* * * * *")
	if err != nil {
		return scheduler.DefaultClient{}, err
	}
	client.Client = &mocks.MockSchedulerClient{}
	return client, nil
}

func TestClient(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "returns scheduler",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				client, err := newMockSchedulerClient()
				g.Expect(err).NotTo(HaveOccurred())
				_ = client.GetScheduler()
			},
		},
		{
			testCase: "returns error for bad cron pattern",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				_, err := scheduler.NewClient("badcron1234")
				g.Expect(err.Error()).To(ContainSubstring("cron expression failed to be parsed"))
			},
		},
		{
			testCase: "returns error for bad load option",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				badFunction := func(client *scheduler.DefaultClient) error {
					return fmt.Errorf("foo")
				}

				_, err := scheduler.NewClient("* * * * *", badFunction)
				g.Expect(err).To(MatchError("foo"))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
