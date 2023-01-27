package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/kopkunka55/appsyncgen/codegen/datasource"
	"github.com/kopkunka55/appsyncgen/codegen/schema"
	"github.com/kopkunka55/appsyncgen/codegen/templates"
	"github.com/kopkunka55/appsyncgen/codegen/templates/resolver"
	"github.com/kopkunka55/appsyncgen/codegen/templates/resolver/dynamodb"
	"github.com/kopkunka55/appsyncgen/codegen/utils"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Operation string

type ResolverList []*Resolver

func (r ResolverList) ForName(name string) *Resolver {
	for _, it := range r {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type Resolver struct {
	*schema.Schema `json:"-"`
	Name           string                `json:"name"`
	Datasource     datasource.DataSource `json:"datasource"`
	Object         schema.Object         `json:"object"`
	ReturnType     schema.ReturnType     `json:"return_type"`
	Functions      []string              `json:"functions"`
}

func NewResolver(field schema.Field, object schema.Object, datasource datasource.DataSource, schema schema.Schema) *Resolver {
	return &Resolver{
		Object:     object,
		Name:       field.Name,
		ReturnType: field.ReturnType,
		Datasource: datasource,
		Functions:  *new([]string),
		Schema:     &schema,
	}
}

func FromJson(pathToJson string) ResolverList {
	var resolvers struct{ Resolvers ResolverList }
	b, err := os.ReadFile(pathToJson)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(b, &resolvers)
	return resolvers.Resolvers
}

func (r *Resolver) createResolverFile(exportPath string, objectName string) *os.File {
	filename := fmt.Sprintf("%d.js", len(r.Functions)+1)
	path := filepath.Join(exportPath, "functions", objectName, r.Name)
	file, err := utils.CreateFile(path, filename)
	if err != nil {
		log.Fatalf("failed to create %s resolver", r.Name)
	}
	r.Functions = append(r.Functions, filepath.Join(path, filename))
	return file
}

func (r *Resolver) GenerateSubscriptionResolver(exportPath string, tmpl *template.Template) {
	file := r.createResolverFile(exportPath, "Subscription")
	templates.ExecuteTemplate(nil, resolver.Subscription, file, tmpl)
}

func (r *Resolver) GenerateFieldResolver(exportPath string, tmpl *template.Template) {
	originalObj := r.Schema.Objects.ForName(r.ReturnType.Name)
	items := originalObj.Fields.ForName("items")
	nextToken := originalObj.Fields.ForName("nextToken")
	if items != nil && nextToken != nil {
		file := r.createResolverFile(exportPath, r.Object.Name)
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK: r.Object.Name,
			SK: items.ReturnType.Name,
		}, string(dynamodb.DdbFieldQuery), file, tmpl)
		file = r.createResolverFile(exportPath, r.Object.Name)
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:        items.ReturnType.Name,
			SK:        items.ReturnType.Name,
			TableName: r.Datasource.Name,
		}, dynamodb.BatchGetItem, file, tmpl)
	} else {
		file := r.createResolverFile(exportPath, r.Object.Name)
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK: originalObj.Name,
			SK: originalObj.Name,
		}, dynamodb.DdbFieldGetItem, file, tmpl)
	}
}
