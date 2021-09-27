package main

import (
	"os"

	"github.com/tristan/myrealdocker/log"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime inplementation.
The purpose of this project is to learn how docker works and how to write a docker by ourselvers
Enjoy it, just for fun.`

var logger = log.NewLogger()

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the defailt ACSII formatter
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("Run failed %v", err.Error)
	}
}



