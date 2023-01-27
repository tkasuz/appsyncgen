package api

import (
	"github.com/kopkunka55/appsyncgen/codegen/datasource"
	"github.com/kopkunka55/appsyncgen/codegen/resolver"
	"github.com/kopkunka55/appsyncgen/codegen/schema"
	"github.com/kopkunka55/appsyncgen/codegen/templates"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type AppSyncApiBuilder struct {
	Name        string
	DataSources datasource.DataSourceList
	Schema      *schema.Schema
	Templates   *template.Template
	ExportPath  *string
}

func NewAppSyncApiBuilder() *AppSyncApiBuilder {
	return &AppSyncApiBuilder{
		DataSources: datasource.DataSourceList{},
	}
}

func (a *AppSyncApiBuilder) SetName(name string) {
	a.Name = name
}

func (a *AppSyncApiBuilder) SetExportPath(exportPath string) {
	cur, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	path := filepath.Clean(filepath.Join(cur, exportPath))
	a.ExportPath = &path
}

func (a *AppSyncApiBuilder) AddDataSource(datasourceType string, datasourceName string) {
	datasource := datasource.NewDataSource(datasourceType, datasourceName)
	if datasource == nil {
		log.Fatalln("Failed to add new datasource")
	}
	a.DataSources = append(a.DataSources, datasource)
}

func (a *AppSyncApiBuilder) SetSchema(pathToSchema string) {
	if a.ExportPath == nil {
		log.Fatalln("export path should be set before setting schema")
	}
	if a.Templates == nil {
		log.Fatalln("export path should be set before setting template")
	}
	cur, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	path := filepath.Clean(filepath.Join(cur, pathToSchema))
	schema := schema.NewSchema(path)
	schema.GenerateNewSchemaFile(*a.ExportPath, a.Templates)
	a.Schema = schema
}

func (a *AppSyncApiBuilder) SetTemplates(pathToTemplate string) {
	tmpl := templates.NewTemplate()
	tmpl = templates.AddFunctionMap(tmpl, map[string]any{
		"toLowerCase": strings.ToLower,
	})
	tmpl = templates.ImportTemplate(tmpl, filepath.Join(pathToTemplate, "resolver", "dynamodb", "*.tmpl"))
	tmpl = templates.ImportTemplate(tmpl, filepath.Join(pathToTemplate, "resolver", "*.tmpl"))
	tmpl = templates.ImportTemplate(tmpl, filepath.Join(pathToTemplate, "graphql", "*.tmpl"))
	a.Templates = tmpl
}

func (a *AppSyncApiBuilder) Build() AppSyncApi {
	if a.Schema == nil {
		log.Fatalln("schema should be set before building AppSync API")
	}
	newSchema := schema.NewSchema(filepath.Join(*a.ExportPath, "schema.graphql"))
	for i, obj := range newSchema.Objects {
		if preObj := a.Schema.Objects.ForName(obj.Name); preObj != nil {
			newSchema.Objects[i].AuthRules = preObj.AuthRules
		}
	}
	resolvers := make(resolver.ResolverList, 0)
	for _, obj := range newSchema.Objects {
		for _, f := range obj.Fields {
			r := resolver.NewResolver(*f, *obj, *a.DataSources[0], *newSchema)
			switch obj.Name {
			case "Mutation":
				r.GenerateMutationResolver(*a.ExportPath, a.Templates)
				resolvers = append(resolvers, r)
			case "Query":
				r.GenerateQueryResolver(*a.ExportPath, a.Templates)
				resolvers = append(resolvers, r)
			case "Subscription":
				r.GenerateSubscriptionResolver(*a.ExportPath, a.Templates)
				resolvers = append(resolvers, r)
			default:
				if a.Schema.Objects.ForName(obj.Name) != nil || strings.HasSuffix(obj.Name, "Payload") {
					if r.ReturnType.IsPrimitive() == false {
						r.GenerateFieldResolver(*a.ExportPath, a.Templates)
						resolvers = append(resolvers, r)
					}
				}
			}
		}
	}
	return AppSyncApi{
		Name:        a.Name,
		DataSources: a.DataSources,
		Schema:      *newSchema,
		Templates:   *a.Templates,
		ExportPath:  *a.ExportPath,
		Resolvers:   resolvers,
	}
}
