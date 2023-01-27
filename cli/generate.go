package cli

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kopkunka55/appsyncgen/codegen/api"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"time"
)

func GenerateAction(context *cli.Context) error {
	s := spinner.New(spinner.CharSets[35], 2*time.Second)
	c := color.New(color.FgHiMagenta)
	s.Color("green", "bold")
	apiBuilder := api.NewAppSyncApiBuilder()
	apiBuilder.SetName(context.String("name"))
	apiBuilder.SetExportPath(context.String("output"))
	apiBuilder.SetTemplates("./codegen/templates")
	apiBuilder.SetSchema(context.String("schema"))
	fmt.Printf("✅ Generated GraphQL schema to %s\n", c.Sprint(filepath.Join(*apiBuilder.ExportPath, "schema.graphql")))
	apiBuilder.AddDataSource("DYNAMODB", context.String("name"))
	appsync := apiBuilder.Build()
	appsync.Export()
	fmt.Printf("✅ Successfully generated JavaScript Resolvers to %s\n", c.Sprint(filepath.Join(*apiBuilder.ExportPath, "functions")))
	fmt.Println("Synthesizing CloudFormation template.....")
	s.Start()
	appsync.Synth()
	s.Stop()
	fmt.Printf("✅ Synthesized CloudFormation Templates to %s\n", c.Sprint(filepath.Join(*apiBuilder.ExportPath, "cloudformation")))
	return nil
}

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen"},
		Usage:   "Generate and prints the AppSync JavaScript Resolver for the given schema [aliases: gen]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "schema",
				Aliases:  []string{"s"},
				Value:    "schema.graphql",
				Usage:    "Path to schema.graphql to generate AppSync JavaScript Resolvers (default: schema.graphql)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Value:    "./build",
				Usage:    "Emits the generated GraphQL schema & resolvers into a directory (default: ./build)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "A user-supplied name for the GraphQL API",
				Required: true,
			},
		},
		Action: GenerateAction,
	}
}
