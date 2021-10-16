package main

import (
	"os"

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

func main() {
	app := controllers.NewQrkDNSApp(Version, Commands)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
