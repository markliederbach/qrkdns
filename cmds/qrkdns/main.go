package main

import (
	"fmt"

	"github.com/markliederbach/qrkdns/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("token: %v", conf.CloudFlareAPIToken)
}
