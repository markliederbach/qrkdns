package config

import (
	"github.com/markliederbach/configr"
	log "github.com/sirupsen/logrus"
)

// Config holds the app configuration
type Config struct {

	// NetworkID is a unique identifier used as the subdomain for
	// this network.
	NetworkID string `env:"NETWORK_ID"`

	// DomainName sets the target zone's domain name
	DomainName string `env:"DOMAIN_NAME" default:"qrkdns.net"`

	// CloudFlareAccountID is the target cloudflare account
	CloudFlareAccountID string `env:"CLOUDFLARE_ACCOUNT_ID"`

	// CloudFlareAPIToken contains the secret scoped API token
	CloudFlareAPIToken string `env:"CLOUDFLARE_API_TOKEN"`

	// LogLevel controls the output verbosity
	LogLevel string `env:"LOG_LEVEL" default:"INFO"`
}

// Load ingests configurations from various sources into a struct
func Load() (Config, error) {

	log.SetFormatter(&log.JSONFormatter{})

	conf := Config{}
	if err := configr.Load(&conf); err != nil {
		return Config{}, err
	}

	logrusLevel, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		return Config{}, err
	}

	log.SetLevel(logrusLevel)

	return conf, nil
}
