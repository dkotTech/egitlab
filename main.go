package main

import (
	"log"
	"os"

	"github.com/dkotTech/egitlab/comands"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:           "egitlab",
		Usage:          "cli extensions for gitlab",
		DefaultCommand: "pipelines",
		Commands: []*cli.Command{
			comands.NewSetCredentialsCommand(),
			comands.NewPipelinesCommand(),
			comands.NewStylesTestCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
