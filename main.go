package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rambler"
	app.Usage = "migrate all the things!"
	app.Version = "2-dev"
	app.Author = "Romain Baugue"
	app.Email = "romain.baugue@elwinar.com"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "configuration, c",
			Value: os.Args[0] + ".json",
			Usage: "path to the configuration file",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "print debug output",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress info output",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "default",
			Usage: "set the working environment",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "apply",
			Usage:  "apply the next migration",
			Action: Apply,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Apply all migrations",
				},
			},
		},
		{
			Name:   "reverse",
			Usage:  "reverse the last migration",
			Action: Reverse,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Reverse all migrations",
				},
			},
		},
	}
	app.Run(os.Args)
}
