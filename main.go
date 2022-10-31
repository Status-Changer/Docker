package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `Docker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		execCommand,
		logCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
