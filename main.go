package main

import (
	"appsyncgen/cli"
	urfave "github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {
	app := &urfave.App{
		Name:                 "appsyncgen",
		Usage:                "Auto generate AppSync JavaScript Resolver",
		Version:              "v1.0.0",
		EnableBashCompletion: true,
		Commands: []*urfave.Command{
			cli.GenerateCommand(),
		},
	}

	sort.Sort(urfave.FlagsByName(app.Flags))
	sort.Sort(urfave.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
