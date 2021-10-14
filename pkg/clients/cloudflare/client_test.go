package cloudflare_test

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cloudflare/cloudflare-go"
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

				records, err := client.ListDNSARecords(ctx, "bar")
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

				_, err = client.ListDNSARecords(ctx, "bar")
				g.Expect(err).To(MatchError("nope"))
			},
		},
		{
			testCase: "apply creates a new record",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				err := configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{},
				)
				g.Expect(err).NotTo(HaveOccurred())

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				expectedRecord := client.BuildDNSARecord("foo", "9.9.9.9")

				err = configrmocks.AddObjectReturns(
					"CreateDNSRecord",
					&sdk.DNSRecordResponse{
						Result: expectedRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				record, err := client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(record).To(Equal(expectedRecord))
			},
		},
		{
			testCase: "apply updates existing record and deletes others",
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

				updateRecord := client.BuildDNSARecord("bar", "1.2.3.4")
				updateRecord.TTL = 9
				updateRecord.ID = "foo"
				deleteRecord := client.BuildDNSARecord("bar", "4.3.2.1")

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{
						updateRecord.ToCloudFlareDNSRecord(),
						deleteRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddObjectReturns(
					"DNSRecord",
					updateRecord.ToCloudFlareDNSRecord(),
				)
				g.Expect(err).NotTo(HaveOccurred())

				record, err := client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).NotTo(HaveOccurred())

				// Mock won't actually update IP, so we just
				// expect the mocked value we passed in
				g.Expect(record).To(Equal(updateRecord))
			},
		},
		{
			testCase: "apply deletes existing and leaves correct record in place",
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

				equalRecord := client.BuildDNSARecord("bar", "1.2.3.4")
				deleteRecord := client.BuildDNSARecord("bar", "5.5.5.5")

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{
						equalRecord.ToCloudFlareDNSRecord(),
						deleteRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				record, err := client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(record).To(Equal(equalRecord))
			},
		},
		{
			testCase: "reports error updating a new record",
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

				updateRecord := client.BuildDNSARecord("bar", "1.2.3.4")
				updateRecord.TTL = 4 // set to something else

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{
						updateRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddErrorReturns(
					"UpdateDNSRecord",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).To(MatchError("baz"))
			},
		},
		{
			testCase: "reports error retrieving newly updated record",
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

				updateRecord := client.BuildDNSARecord("bar", "1.2.3.4")
				updateRecord.TTL = 4 // set to something else

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{
						updateRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddErrorReturns(
					"DNSRecord",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).To(MatchError("baz"))
			},
		},
		{
			testCase: "reports error deleting while creating a new record",
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

				existingRecord := client.BuildDNSARecord("bar", "5.5.5.5")
				deleteRecord := client.BuildDNSARecord("bar", "5.5.5.6")

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{
						existingRecord.ToCloudFlareDNSRecord(),
						deleteRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddErrorReturns(
					"DeleteDNSRecord",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).To(MatchError("baz"))
			},
		},
		{
			testCase: "reports error creating a new record",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				err := configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{},
				)
				g.Expect(err).NotTo(HaveOccurred())

				err = configrmocks.AddErrorReturns(
					"CreateDNSRecord",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).To(MatchError("baz"))
			},
		},
		{
			testCase: "reports error listing existing records",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				ctx := context.Background()

				err := configrmocks.AddErrorReturns(
					"DNSRecords",
					fmt.Errorf("boo"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				client, err := cloudflare.NewClientWithToken(
					ctx,
					"account1234",
					"foo.net",
					"token1234",
					withMockSDKClient,
				)
				g.Expect(err).NotTo(HaveOccurred())

				_, err = client.ApplyDNSARecord(ctx, "bar", "1.2.3.4")
				g.Expect(err).To(MatchError("boo"))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
