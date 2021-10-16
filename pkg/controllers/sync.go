package controllers

import (
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	// CloudflareClientOptions is used by testing to inject a mock client option
	CloudflareClientOptions = []cloudflare.LoadOption{}
	// IPClientOptions is used by testing to inject a mock client option
	IPClientOptions = []ip.LoadOption{}
)

const (
	// NetworkIDFlag wraps the name of the command flag
	NetworkIDFlag string = "network-id"

	// DomainFlag wraps the name of the command flag
	DomainFlag string = "domain"

	// CloudflareAccountIDFlag wraps the name of the command flag
	CloudflareAccountIDFlag string = "cf-account-id"

	// CloudflareAPITokenFlag wraps the name of the command flag
	CloudflareAPITokenFlag string = "cf-api-token"

	// IPServiceURLFlag wraps the name of the command flag
	IPServiceURLFlag string = "ip-service-url"
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
				Name:    DomainFlag,
				Aliases: []string{"d"},
				Usage:   "Base domain used when constructing the host's subdomain",
				EnvVars: []string{"DOMAIN_NAME"},
				Value:   "qrkdns.net",
			},
			&cli.StringFlag{
				Name:     CloudflareAccountIDFlag,
				Aliases:  []string{"a"},
				Usage:    "Cloudflare Account ID used in conjunction with the API token",
				EnvVars:  []string{"CLOUDFLARE_ACCOUNT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     CloudflareAPITokenFlag,
				Aliases:  []string{"t"},
				Usage:    "Cloudflare API token providing scoped permisions for DNS management",
				EnvVars:  []string{"CLOUDFLARE_API_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    IPServiceURLFlag,
				Aliases: []string{"i"},
				Usage:   "Web service to retrieve external IP address",
				EnvVars: []string{"IP_SERVICE_URL"},
				Value:   "http://checkip.amazonaws.com",
			},
		},
		Action: syncOnce,
	}
}

func syncOnce(c *cli.Context) error {
	accountID := c.String(CloudflareAccountIDFlag)
	domain := c.String(DomainFlag)
	apiToken := c.String(CloudflareAPITokenFlag)
	ipServiceURL := c.String(IPServiceURLFlag)
	networkID := c.String(NetworkIDFlag)

	cloudflareClient, err := cloudflare.NewClientWithToken(
		c.Context,
		accountID,
		domain,
		apiToken,
		CloudflareClientOptions...,
	)
	if err != nil {
		return err
	}

	ipClient, err := ip.NewClient(ipServiceURL, IPClientOptions...)
	if err != nil {
		return err
	}

	externalIP, err := ipClient.GetExternalIPAddress(c.Context)
	if err != nil {
		return err
	}

	_, err = cloudflareClient.ApplyDNSARecord(c.Context, networkID, externalIP)
	if err != nil {
		return err
	}

	log.Infof("Sync complete")
	return nil
}
