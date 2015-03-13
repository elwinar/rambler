package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
)

var app *cli.App
var service *Service

func init() {
	log.SetFlags(0)

	app = cli.NewApp()
	app.Name = "rambler"
	app.Usage = "Migrate all the things!"
	app.Version = "3"
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
	}
	app.Before = Bootstrap
	app.Commands = []cli.Command{
		{
			Name:   "apply",
			Usage:  "apply the next migration",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Apply all migrations",
				},
			},
			Action: Apply,
		},
		// 		{
		// 			Name:   "reverse",
		// 			Usage:  "reverse the last migration",
		// 			Flags: []cli.Flag{
		// 				cli.BoolFlag{
		// 					Name:  "all, a",
		// 					Usage: "Reverse all migrations",
		// 				},
		// 			},
		// 			Action: Reverse,
		// 		},
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}
