package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/dns"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
	"github.com/markliederbach/qrkdns/pkg/clients/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	// CloudflareClientOptions is used by testing to inject a mock client option
	CloudflareClientOptions = []cloudflare.LoadOption{}
	// IPClientOptions is used by testing to inject a mock client option
	IPClientOptions = []ip.LoadOption{}
	// SchedulerClientOptions is used by testing to inject a mock client option
	SchedulerClientOptions = []scheduler.LoadOption{}
)

const (
	// NetworkIDFlag wraps the name of the command flag
	NetworkIDFlag string = "network-id"

	// DomainFlag wraps the name of the command flag
	DomainFlag string = "domain"

	// ProviderTypeFlag wraps the name of the command flag
	ProviderTypeFlag string = "provider"

	// CloudflareAccountIDFlag wraps the name of the command flag
	CloudflareAccountIDFlag string = "cf-account-id"

	// CloudflareAPITokenFlag wraps the name of the command flag
	CloudflareAPITokenFlag string = "cf-api-token"

	// IPServiceURLFlag wraps the name of the command flag
	IPServiceURLFlag string = "ip-service-url"

	// TimeoutFlag wraps the name of the command flag
	TimeoutFlag string = "timeout"

	// ScheduleFlag wraps the name of the command flag
	ScheduleFlag string = "schedule"
)

// SyncCommand returns
func SyncCommand() *cli.Command {
	return &cli.Command{
		Name:    "sync",
		Aliases: []string{"s"},
		Usage:   "Sync this host's external IP to Cloudflare",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     NetworkIDFlag,
				Aliases:  []string{"n"},
				Usage:    "Identifier used for the subdomain",
				EnvVars:  []string{"NETWORK_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     DomainFlag,
				Aliases:  []string{"d"},
				Usage:    "Base domain used when constructing the host's subdomain",
				EnvVars:  []string{"DOMAIN_NAME"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    ProviderTypeFlag,
				Aliases: []string{"p"},
				Usage:   fmt.Sprintf("Type of provider to use (one of: %v)", getSupportedProvidersString()),
				EnvVars: []string{"PROVIDER"},
				Value:   string(dns.ProviderTypeCloudflare),
			},
			&cli.StringFlag{
				Name:    CloudflareAccountIDFlag,
				Aliases: []string{"a"},
				Usage:   "Cloudflare Account ID used in conjunction with the API token",
				EnvVars: []string{"CLOUDFLARE_ACCOUNT_ID"},
			},
			&cli.StringFlag{
				Name:    CloudflareAPITokenFlag,
				Aliases: []string{"t"},
				Usage:   "Cloudflare API token providing scoped permisions for DNS management",
				EnvVars: []string{"CLOUDFLARE_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    IPServiceURLFlag,
				Aliases: []string{"i"},
				Usage:   "Web service to retrieve external IP address",
				EnvVars: []string{"IP_SERVICE_URL"},
				Value:   "http://checkip.amazonaws.com",
			},
			&cli.StringFlag{
				Name:    TimeoutFlag,
				Usage:   "Timeout as a duration string (e.g., 5s). Empty/Unset means no timeout",
				Value:   "",
				EnvVars: []string{"TIMEOUT"},
			},
		},
		Action: syncOnce,
		Subcommands: []*cli.Command{
			{
				Name:  "cron",
				Usage: "Run the sync on a recurring schedule",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     ScheduleFlag,
						Usage:    "Cron pattern",
						EnvVars:  []string{"SCHEDULE"},
						Required: true,
					},
				},
				Action: syncCron,
			},
		},
	}
}

// syncOnce performs a single sync task. Each sync consists of
// retrieving the external IP Address of this host and applying
// the result as a DNS A record to the specified provider through a
// dedicated API client.
func syncOnce(c *cli.Context) error {
	var cancel context.CancelFunc

	ctx := c.Context
	ipServiceURL := c.String(IPServiceURLFlag)
	networkID := c.String(NetworkIDFlag)
	timeoutString := c.String(TimeoutFlag)

	if timeoutString != "" {
		timeoutDuration, err := time.ParseDuration(timeoutString)
		if err != nil {
			log.WithError(err).Error("Failed to parse timeout duration")
			return err
		}
		log.WithField("timeout", timeoutDuration).Debug("Setting timeout")
		ctx, cancel = context.WithTimeout(c.Context, timeoutDuration)
		defer cancel()
	}

	dnsClient, err := buildDNSProvider(c)
	if err != nil {
		log.WithError(err).Error("Failed to build DNS client")
		return err
	}

	ipClient, err := ip.NewClient(ipServiceURL, IPClientOptions...)
	if err != nil {
		log.WithError(err).Error("Failed to build IP client")
		return err
	}

	externalIP, err := ipClient.GetExternalIPAddress(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to get external IP address")
		return err
	}

	log.WithField("externalIP", externalIP).Debug("External IP address retrieved")

	_, err = dnsClient.ApplyDNSARecord(ctx, networkID, externalIP)
	if err != nil {
		log.WithError(err).Error("Failed to apply DNS A record")
		return err
	}

	log.Info("Sync complete")
	return nil
}

// syncCron runs the syncOnce task at the specified cron schedule
func syncCron(c *cli.Context) error {
	scheduleCron := c.String(ScheduleFlag)

	cronLog := log.WithField("schedule", scheduleCron)

	client, err := scheduler.NewClient(scheduleCron, SchedulerClientOptions...)
	if err != nil {
		return err
	}

	clientScheduler := client.GetScheduler()

	_, err = clientScheduler.Do(syncOnce, c)
	if err != nil {
		return err
	}

	cronLog.Info("Running cron scheduler")
	clientScheduler.StartBlocking() // does not return

	return nil
}

// buildDNSProvider determines which provider to create and returns
// an instantiated provider client
func buildDNSProvider(c *cli.Context) (dns.Provider, error) {
	var dnsClient dns.Provider
	var err error

	ctx := c.Context
	providerType := c.String(ProviderTypeFlag)

	switch dns.ProviderType(providerType) {
	case dns.ProviderTypeCloudflare:
		var cloudflareOptions map[string]string

		domain := c.String(DomainFlag)
		cloudflareOptions, err = stringsOrError(
			c,
			fmt.Sprintf("using %s provider", dns.ProviderTypeCloudflare),
			CloudflareAccountIDFlag,
			CloudflareAPITokenFlag,
		)
		if err != nil {
			return dnsClient, err
		}

		dnsClient, err = cloudflare.NewClientWithToken(
			ctx,
			cloudflareOptions[CloudflareAccountIDFlag],
			domain,
			cloudflareOptions[CloudflareAPITokenFlag],
			CloudflareClientOptions...,
		)
		if err != nil {
			return dnsClient, err
		}
	default:
		return dnsClient, fmt.Errorf("unsupported DNS client: %v", providerType)
	}
	return dnsClient, nil
}

// getSupportedProvidersString returns the supported provider types
// as a comma-separated string
func getSupportedProvidersString() string {
	stringProviders := []string{}
	for _, provider := range dns.SuportedProviders {
		stringProviders = append(stringProviders, string(provider))
	}
	return strings.Join(stringProviders, ", ")
}

// stringsOrError attempts to load a list of options from the CLI, and reports
// back any missing presumably required options
func stringsOrError(c *cli.Context, whenMessage string, options ...string) (map[string]string, error) {
	results := make(map[string]string)
	missingOptions := []string{}
	for _, option := range options {
		value := c.String(option)
		if value == "" {
			missingOptions = append(missingOptions, option)
			continue
		}
		results[option] = value
	}
	if len(missingOptions) > 0 {
		formattedOptions := []string{}
		for _, option := range options {
			formattedOptions = append(formattedOptions, fmt.Sprintf("--%v", option))
		}
		return make(map[string]string), fmt.Errorf(
			"options [%v] are required when %v",
			strings.Join(formattedOptions, ", "),
			whenMessage,
		)
	}
	return results, nil
}
