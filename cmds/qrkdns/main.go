package main

import (
	"context"

	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
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

	_, err = cloudflareClient.ListDNSRecords(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ipClient := ip.NewClient(conf.IPServiceURL)
	externalIP, err := ipClient.GetExternalIPAddress(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("External IP Address: %v", externalIP)
}
