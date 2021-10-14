package ip_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	configrmocks "github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
	"github.com/markliederbach/qrkdns/pkg/mocks"
	. "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func newMockIPClient() (ip.DefaultClient, error) {
	client, err := ip.NewClient("some_url")
	if err != nil {
		return ip.DefaultClient{}, err
	}
	client.Client = &mocks.MockHTTPClient{}
	return client, nil
}

func TestFile(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "returns ip address",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				ctx := context.Background()
				client, err := newMockIPClient()
				g.Expect(err).NotTo(HaveOccurred())

				ipAddress, err := client.GetExternalIPAddress(ctx)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(ipAddress).To(Equal(mocks.DefaultExternalIPAddress))
			},
		},
		{
			testCase: "returns error from external ip address lookup",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				ctx := context.Background()
				client, err := newMockIPClient()
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddErrorReturns("Get", fmt.Errorf("oh no"))
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.GetExternalIPAddress(ctx)
				g.Expect(err).To(MatchError("oh no"))
			},
		},
		{
			testCase: "returns error from external ip address lookup body reader",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				ctx := context.Background()
				client, err := newMockIPClient()
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddObjectReturns(
					"Get",
					&http.Response{
						StatusCode: 200,
						Body:       &mocks.ErrorReader{Error: fmt.Errorf("error reader")},
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.GetExternalIPAddress(ctx)
				g.Expect(err).To(MatchError("error reader"))
			},
		},
		{
			testCase: "returns error from external ip address lookup body reader",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				ctx := context.Background()
				client, err := newMockIPClient()
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddObjectReturns(
					"Get",
					&http.Response{
						StatusCode: 200,
						Body:       &mocks.ErrorReader{Error: fmt.Errorf("error reader")},
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.GetExternalIPAddress(ctx)
				g.Expect(err).To(MatchError("error reader"))
			},
		},
		{
			testCase: "returns error from external ip address lookup status code",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				ctx := context.Background()
				client, err := newMockIPClient()
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddObjectReturns(
					"Get",
					&http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(strings.NewReader("foo")),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.GetExternalIPAddress(ctx)
				g.Expect(err).To(HaveOccurred())
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
