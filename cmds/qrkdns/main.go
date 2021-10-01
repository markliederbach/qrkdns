package main

import (
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
	"github.com/markliederbach/qrkdns/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	cloudflareClient, err := cloudflare.NewCloudflareClient(
		conf.CloudFlareAccountID,
		conf.DomainName,
		conf.CloudFlareAPIToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Domain %v has zone ID %v", cloudflareClient.DomainName, cloudflareClient.ZoneID)
}
