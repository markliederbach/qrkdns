package main

import (
	"os"
	"time"

	"github.com/markliederbach/qrkdns/pkg/controllers"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var (
	// Version tracks the semantic version of this release
	Version = "latest"
	// Commands contains the base commands to attach to this CLI
	Commands = []*cli.Command{
		controllers.SyncCommand(),
	}
)

const (
	// LogLevelFlag wraps the name of the command flag
	LogLevelFlag string = "log-level"
)

// NewApp creates a new CLI app
func NewApp(version string, commands []*cli.Command) *cli.App {
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
			log.WithField("version", Version).Debug("Running qrkdns")
			return nil
		},
		Commands: commands,
	}
}

func main() {
	app := NewApp(Version, Commands)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
