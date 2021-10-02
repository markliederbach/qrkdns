package main

import (
	"context"

	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.TODO()
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	cloudflareClient, err := cloudflare.NewClientWithToken(
		ctx,
		conf.CloudFlareAccountID,
		conf.DomainName,
		conf.CloudFlareAPIToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	dnsRecords, err := cloudflareClient.ListDNSRecords(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Domain %v has %v DNS records", cloudflareClient.DomainName, len(dnsRecords))
}
