package controllers_test

import (
	"fmt"
	"testing"

	sdk "github.com/cloudflare/cloudflare-go"
	configrmocks "github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
	"github.com/markliederbach/qrkdns/pkg/clients/scheduler"
	"github.com/markliederbach/qrkdns/pkg/controllers"
	"github.com/markliederbach/qrkdns/pkg/mocks"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
)

func withMockSDKClient(client *cloudflare.DefaultClient) error {
	client.Client = &mocks.MockCloudflareSDKClient{}
	return nil
}

func withMockHTTPClient(client *ip.DefaultClient) error {
	client.Client = &mocks.MockHTTPClient{}
	return nil
}

func withMockSchedulerClient(client *scheduler.DefaultClient) error {
	client.Client = &mocks.MockSchedulerClient{}
	return nil
}

func TestSync(t *testing.T) {
	controllers.CloudflareClientOptions = append(
		controllers.CloudflareClientOptions,
		withMockSDKClient,
	)
	controllers.IPClientOptions = append(
		controllers.IPClientOptions,
		withMockHTTPClient,
	)

	// disable help text for tests
	cli.AppHelpTemplate = ""

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
						"TIMEOUT":               "1s",
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

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			testCase: "returns error for bad timeout string",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
						"TIMEOUT":               "badtimeout1234",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(MatchError("time: invalid duration \"badtimeout1234\""))
			},
		},
		{
			testCase: "returns error for missing required option",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "", // missing required flag
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(HaveOccurred())
			},
		},
		{
			testCase: "returns error for new cloudflare client error",
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

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(HaveOccurred())
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

				oldIPClientOptions := controllers.IPClientOptions
				defer func() {
					controllers.IPClientOptions = oldIPClientOptions
				}()

				controllers.IPClientOptions = append(
					controllers.IPClientOptions,
					func(client *ip.DefaultClient) error {
						return fmt.Errorf("boo")
					},
				)

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(MatchError("boo"))
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
					"Do",
					fmt.Errorf("baz"),
				)
				g.Expect(err).NotTo(HaveOccurred())

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(MatchError("baz"))
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

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync"})
				g.Expect(err).To(MatchError("baz"))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}

func TestSyncCron(t *testing.T) {
	controllers.CloudflareClientOptions = append(
		controllers.CloudflareClientOptions,
		withMockSDKClient,
	)
	controllers.IPClientOptions = append(
		controllers.IPClientOptions,
		withMockHTTPClient,
	)
	controllers.SchedulerClientOptions = append(
		controllers.SchedulerClientOptions,
		withMockSchedulerClient,
	)

	// disable help text for tests
	cli.AppHelpTemplate = ""

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
						"TIMEOUT":               "1s",
						"SCHEDULE":              "* * * * *",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync", "cron"})
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			testCase: "returns error from new scheduler client",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
						"TIMEOUT":               "1s",
						"SCHEDULE":              "badcron1234",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync", "cron"})
				g.Expect(err.Error()).To(ContainSubstring("cron expression failed to be parsed"))
			},
		},
		{
			testCase: "returns error from new scheduler client",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := configrmocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"NETWORK_ID":            "xxx",
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
						"TIMEOUT":               "1s",
						"SCHEDULE":              "* * * * *",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				err = configrmocks.AddErrorReturns("Do", fmt.Errorf("foo"))
				g.Expect(err).NotTo(HaveOccurred())

				app := controllers.NewQrkDNSApp(
					"version123",
					[]*cli.Command{controllers.SyncCommand()},
				)

				err = app.Run([]string{"qrkdns", "sync", "cron"})
				g.Expect(err).To(MatchError("foo"))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
