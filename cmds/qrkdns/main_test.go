package main

import (
	"fmt"
	"testing"

	sdk "github.com/cloudflare/cloudflare-go"
	configrmocks "github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
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

func withMockHTTPClient(client *ip.DefaultClient) error {
	client.Client = &mocks.MockHTTPClient{}
	return nil
}

func TestMain(t *testing.T) {
	CloudflareClientOptions = append(
		CloudflareClientOptions,
		withMockSDKClient,
	)
	IPClientOptions = append(
		IPClientOptions,
		withMockHTTPClient,
	)

	tests := []testRunner{
		{
			testCase: "runs successfully",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				err = configrmocks.AddObjectReturns(
					"DNSRecords",
					[]sdk.DNSRecord{},
				)
				g.Expect(err).NotTo(HaveOccurred())

				expectedRecord := cloudflare.BuildDNSARecord("foo", "foo.net", mocks.DefaultExternalIPAddress)

				err = configrmocks.AddObjectReturns(
					"CreateDNSRecord",
					&sdk.DNSRecordResponse{
						Result: expectedRecord.ToCloudFlareDNSRecord(),
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(main).NotTo(Panic())
			},
		},
		{
			testCase: "panics for config error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						// "NETWORK_ID":            "xxx", // missing required variable
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				g.Expect(main).To(Panic())
			},
		},
		{
			testCase: "panics for new cloudflare client error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				err = configrmocks.AddErrorReturns(
					"ZoneIDByName",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(main).To(Panic())
			},
		},
		{
			testCase: "panics for new ip client error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				oldIPClientOptions := IPClientOptions
				defer func() {
					IPClientOptions = oldIPClientOptions
				}()

				IPClientOptions = append(
					IPClientOptions,
					func(client *ip.DefaultClient) error {
						return fmt.Errorf("boo")
					},
				)

				g.Expect(main).To(Panic())
			},
		},
		{
			testCase: "panics for error from getting external ip",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				err = configrmocks.AddErrorReturns(
					"Get",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(main).To(Panic())
			},
		},
		{
			testCase: "panics for error from applying record",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				err = configrmocks.AddErrorReturns(
					"DNSRecords",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(main).To(Panic())
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
