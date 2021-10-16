package controllers

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	// LogLevelFlag wraps the name of the command flag
	LogLevelFlag string = "log-level"
)

// NewQrkDNSApp creates a new CLI app
func NewQrkDNSApp(version string, commands []*cli.Command) *cli.App {
	return &cli.App{
		Name:     "qrkdns",
		Version:  version,
		Usage:    "Automatically update Cloudflare DNS records",
		Compiled: time.Now(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    LogLevelFlag,
				Aliases: []string{"l"},
				Value:   "INFO",
				Usage:   "Set the log output",
				EnvVars: []string{"LOG_LEVEL"},
			},
		},
		Before: func(c *cli.Context) error {
			log.SetFormatter(&log.JSONFormatter{})
			logrusLevel, err := log.ParseLevel(c.String(LogLevelFlag))
			if err != nil {
				return err
			}
			log.SetLevel(logrusLevel)
			log.WithField("version", version).Debug("Running qrkdns")
			return nil
		},
		Commands: commands,
	}
}
