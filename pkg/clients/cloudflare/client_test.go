package cloudflare_test

import (
	"context"
	"fmt"
	"testing"

	configrmocks "github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/mocks"
	. "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func withMockSDKClient(client *cloudflare.DefaultClient) error {
	client.Client = &mocks.MockCloudflareSDKClient{}
	return nil
}

func TestFile(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "lists dns records",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				records, err := client.ListDNSRecords(ctx)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(records).To(Equal(mocks.DefaultDNSRecords))
			},
		},
		{
			testCase: "returns cached zone ID",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				configrmocks.AddObjectReturns("ZoneIDByName", "newzone")

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				// Would otherwise have returned the default mock zone
				zoneID, err := client.GetZoneID(ctx)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(zoneID).To(Equal("newzone"))
			},
		},
		{
			testCase: "load option returns error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				_, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					func(client *cloudflare.DefaultClient) error {
						return fmt.Errorf("oh no")
					},
				)
				g.Expect(err).To(MatchError("oh no"))
			},
		},
		{
			testCase: "returns error from token initializer",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				_, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"", // errors when token is empty
				)
				g.Expect(err).To(HaveOccurred())
			},
		},
		{
			testCase: "returns error from Zone ID preloader",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				configrmocks.AddErrorReturns("ZoneIDByName", fmt.Errorf("no no no"))

				_, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).To(MatchError("no no no"))
			},
		},
		{
			testCase: "returns error from dns records",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				configrmocks.AddErrorReturns("DNSRecords", fmt.Errorf("nope"))

				_, err = client.ListDNSRecords(ctx)
				g.Expect(err).To(MatchError("nope"))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
