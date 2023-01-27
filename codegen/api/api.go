package api

import (
	"appsyncgen/codegen/datasource"
	"appsyncgen/codegen/resolver"
	"appsyncgen/codegen/schema"
	"appsyncgen/codegen/utils"
	"encoding/json"
	"log"
	"os"
	"text/template"
)

type AppSyncApi struct {
	Name        string                    `json:"name"`
	DataSources datasource.DataSourceList `json:"-"`
	Schema      schema.Schema             `json:"schema"`
	Templates   template.Template         `json:"-"`
	ExportPath  string                    `json:"export_path"`
	Resolvers   resolver.ResolverList     `json:"resolvers"`
}

func (a AppSyncApi) Export() {
	data := utils.ToJson(a)
	file, err := utils.CreateFile(a.ExportPath, "resolvers.json")
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := file.WriteString(data); err != nil {
		log.Fatalln(err)
	}
}

func FromJson(pathToJson string) *AppSyncApi {
	b, err := os.ReadFile(pathToJson)
	if err != nil {
		log.Fatalln(err)
	}
	api := AppSyncApi{}
	json.Unmarshal(b, &api)
	return &api
}
