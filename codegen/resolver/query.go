package resolver

import (
	"fmt"
	"github.com/kopkunka55/appsyncgen/codegen/templates"
	"github.com/kopkunka55/appsyncgen/codegen/templates/resolver/dynamodb"
	"strings"
	"text/template"
)

func (r *Resolver) generateGetResolver(exportPath string, tmpl *template.Template) {
	file := r.createResolverFile(exportPath, "Query")
	templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
		PK: r.ReturnType.Name,
		SK: r.ReturnType.Name,
	}, dynamodb.GetItem, file, tmpl)
}

func (r *Resolver) generateListResolver(exportPath string, tmpl *template.Template) {
	returnedType := r.Schema.Objects.ForName(r.ReturnType.Name)
	originalType := returnedType.Fields.ForName("items").ReturnType.Name
	parentObjName := r.Name[len(fmt.Sprintf("list%ssBy", originalType)):len(r.Name)]
	file := r.createResolverFile(exportPath, "Query")
	templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
		PK: parentObjName,
		SK: originalType,
	}, dynamodb.Query, file, tmpl)
	file = r.createResolverFile(exportPath, "Query")
	templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
		PK:        originalType,
		SK:        originalType,
		TableName: r.Datasource.Name,
	}, dynamodb.BatchGetItem, file, tmpl)
}

func (r *Resolver) GenerateQueryResolver(exportPath string, tmpl *template.Template) {
	if fmt.Sprintf("get%s", r.ReturnType.Name) == r.Name {
		r.generateGetResolver(exportPath, tmpl)
	} else if yes := strings.HasPrefix(r.Name, "list"); yes {
		r.generateListResolver(exportPath, tmpl)
	}
}
