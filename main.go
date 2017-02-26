package main

import (
	"fmt"
	"os"

	"github.com/elwinar/rambler/log"
	"github.com/urfave/cli"
)

var service *Service
var logger *log.Logger

// VERSION holds the version of rambler as defined at compile time.
var VERSION string

func main() {
	var app = cli.NewApp()

	app.Name = "rambler"
	app.Usage = "Migrate all the things!"
	app.Version = VERSION
	app.Authors = []cli.Author{
		{
			Name:  "Romain Baugue",
			Email: "romain.baugue@elwinar.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "configuration, c",
			Value: "rambler.json",
			Usage: "path to the configuration file",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "default",
			Usage: "set the working environment",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "display debug messages",
		},
	}

	app.Before = Bootstrap

	app.Commands = []cli.Command{
		{
			Name:  "apply",
			Usage: "apply the next migration",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Apply all migrations",
				},
			},
			Action: Apply,
		},
		{
			Name:  "reverse",
			Usage: "reverse the last migration",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Reverse all migrations",
				},
			},
			Action: Reverse,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
