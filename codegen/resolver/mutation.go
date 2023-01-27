package resolver

import (
	"appsyncgen/codegen/templates"
	"appsyncgen/codegen/templates/resolver/dynamodb"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

func (r *Resolver) generateCreateResolver(exportPath string, tmpl *template.Template) {
	if r.Schema.Objects.ForName(r.ReturnType.Name).Type.IsConnectionPayload() {
		fields := r.Schema.Objects.ForName(r.ReturnType.Name).Fields
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:        fields[1].ReturnType.Name,
			SK:        fields[3].ReturnType.Name,
			TableName: r.Datasource.Name,
		}, dynamodb.PutConnection, file, tmpl)
		return
	}
	returnType := r.Schema.Objects.ForName(r.ReturnType.Name).Fields[1].ReturnType.Name
	fields := r.Schema.Objects.ForName(fmt.Sprintf("Create%sInput", returnType)).Fields
	attributes := fields.Names()
	if fields.HasConnection() {
		connections := []templates.DynamoDBResolverTemplateData{}
		for _, f := range fields {
			if name, ok := f.IsConnection(); ok {
				connections = append(connections, templates.DynamoDBResolverTemplateData{
					PK:        returnType,
					SK:        fmt.Sprintf("%c%s", strings.ToUpper(*name)[0], (*name)[1:]),
					TableName: r.Datasource.Name,
				})
			}
		}
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:          returnType,
			SK:          returnType,
			TableName:   r.Datasource.Name,
			Attributes:  attributes,
			Connections: connections,
		}, dynamodb.TransactPutItems, file, tmpl)
	} else {
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:         returnType,
			SK:         returnType,
			Attributes: attributes,
		}, dynamodb.PutItem, file, tmpl)
	}
}

func (r *Resolver) generateUpdateResolver(exportPath string, tmpl *template.Template) {
	returnType := r.Schema.Objects.ForName(r.ReturnType.Name).Fields[1].ReturnType.Name
	fields := r.Schema.Objects.ForName(fmt.Sprintf("Update%sInput", returnType)).Fields
	obj := r.Schema.Objects.ForName(returnType)
	if obj.Type.Name == "Connection" {
		return
	}
	fileForGet := r.createResolverFile(exportPath, "Mutation")
	templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
		PK: returnType,
		SK: returnType,
	}, dynamodb.GetItem, fileForGet, tmpl)
	attributes := fields.Names()
	file := r.createResolverFile(exportPath, "Mutation")
	templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
		PK:         returnType,
		SK:         returnType,
		Attributes: attributes,
	}, dynamodb.UpdateItem, file, tmpl)
}

func (r *Resolver) generateDeleteResolver(exportPath string, tmpl *template.Template) {
	if r.Schema.Objects.ForName(r.ReturnType.Name).Type.IsConnectionPayload() {
		fields := r.Schema.Objects.ForName(r.ReturnType.Name).Fields
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:        fields[1].ReturnType.Name,
			SK:        fields[3].ReturnType.Name,
			TableName: r.Datasource.Name,
		}, dynamodb.DeleteConnection, file, tmpl)
		return
	}
	returnType := r.Schema.Objects.ForName(r.ReturnType.Name).Fields[1].ReturnType.Name
	fields := r.Schema.Objects.ForName(returnType).Fields
	if fields.HasConnection() {
		fileForGet := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK: returnType,
			SK: returnType,
		}, dynamodb.GetItem, fileForGet, tmpl)
		connections := []templates.DynamoDBResolverTemplateData{}
		for _, f := range fields {
			if name, ok := f.IsConnection(); ok {
				connections = append(connections, templates.DynamoDBResolverTemplateData{
					PK:        returnType,
					SK:        fmt.Sprintf("%c%s", strings.ToUpper(*name)[0], (*name)[1:]),
					TableName: r.Datasource.Name,
				})
			}
		}
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK:          returnType,
			SK:          returnType,
			TableName:   returnType,
			Connections: connections,
		}, dynamodb.TransactDeleteItems, file, tmpl)
	} else {
		file := r.createResolverFile(exportPath, "Mutation")
		templates.ExecuteTemplate(templates.DynamoDBResolverTemplateData{
			PK: returnType,
			SK: returnType,
		}, dynamodb.DeleteItem, file, tmpl)
	}
}

func (r *Resolver) GenerateMutationResolver(exportPath string, tmpl *template.Template) {
	if yes, _ := regexp.MatchString(`create*`, r.Name); yes {
		r.generateCreateResolver(exportPath, tmpl)
	} else if yes, _ := regexp.MatchString(`update*`, r.Name); yes {
		r.generateUpdateResolver(exportPath, tmpl)
	} else if yes, _ := regexp.MatchString(`delete*`, r.Name); yes {
		r.generateDeleteResolver(exportPath, tmpl)
	}
}
