package main

import (
	"context"

	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
	"github.com/markliederbach/qrkdns/pkg/config"
	log "github.com/sirupsen/logrus"
)

var (
	// CloudflareClientOptions is used by testing to inject a mock client option
	CloudflareClientOptions = []cloudflare.LoadOption{}
	// IPClientOptions is used by testing to inject a mock client option
	IPClientOptions = []ip.LoadOption{}
)

func main() {
	ctx := context.TODO()
	conf, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	cloudflareClient, err := cloudflare.NewClientWithToken(
		ctx,
		conf.CloudFlareAccountID,
		conf.DomainName,
		conf.CloudFlareAPIToken,
		CloudflareClientOptions...,
	)
	if err != nil {
		log.Panic(err)
	}

	ipClient, err := ip.NewClient(conf.IPServiceURL, IPClientOptions...)
	if err != nil {
		log.Panic(err)
	}

	externalIP, err := ipClient.GetExternalIPAddress(ctx)
	if err != nil {
		log.Panic(err)
	}

	_, err = cloudflareClient.ApplyDNSARecord(ctx, conf.NetworkID, externalIP)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("Sync complete")
}
