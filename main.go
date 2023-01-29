package main

import (
	"log"
	"os"
	"sort"

	"github.com/kopkunka55/appsyncgen/cli"
	urfave "github.com/urfave/cli/v2"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	app := &urfave.App{
		Name:                 "appsyncgen",
		Usage:                "Auto generate AppSync JavaScript Resolver",
		Version:              Version,
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
