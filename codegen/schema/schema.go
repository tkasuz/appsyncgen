package schema

import (
	"appsyncgen/codegen/directive"
	"appsyncgen/codegen/templates"
	"appsyncgen/codegen/templates/graphql"
	"appsyncgen/codegen/utils"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/parser"
	"log"
	"text/template"
)

type Schema struct {
	Objects     ObjectList     `json:"-"`
	Enums       EnumList       `json:"-"`
	FilePath    string         `json:"file_path"`
	Connections utils.PairList `json:"-"`
	HasManyList utils.PairList `json:"-"`
}

func (s *Schema) GenerateNewSchemaFile(exportPath string, tmpl *template.Template) {
	f, err := utils.CreateFile(exportPath, "schema.graphql")
	if err != nil {
		log.Fatalln("Unable to generate schema, given export path might be wrong", err)
	}
	templates.ExecuteTemplate(s.Objects, graphql.Type, f, tmpl)
	templates.ExecuteTemplate(s, graphql.ConnectionType, f, tmpl)
	templates.ExecuteTemplate(s.Objects, graphql.Payload, f, tmpl)
	templates.ExecuteTemplate(s, graphql.ConnectionPayload, f, tmpl)
	templates.ExecuteTemplate(s.Objects, graphql.ListType, f, tmpl)
	templates.ExecuteTemplate(s, graphql.Mutation, f, tmpl)
	templates.ExecuteTemplate(s.Objects, graphql.Query, f, tmpl)
	templates.ExecuteTemplate(s.Objects, graphql.Subscription, f, tmpl)
	templates.ExecuteTemplate(struct {
		Objects     ObjectList
		HasManyList utils.PairList
	}{s.Objects, s.HasManyList}, graphql.Input, f, tmpl)
	templates.ExecuteTemplate(s.Connections, graphql.ConnectionInput, f, tmpl)
	templates.ExecuteTemplate(s.Enums, graphql.Enum, f, tmpl)
	templates.ExecuteTemplate(s, graphql.SubscriptionFilter, f, tmpl)
	templates.ExecuteTemplate(nil, graphql.SubscriptionInput, f, tmpl)
	utils.FormatGraphqlSchema(f)
}

func NewSchema(importPath string) *Schema {
	s, err := utils.ReadFile(importPath)
	if err != nil {
		log.Fatalln("Unable to import schema, given import path might be wrong", err)
	}
	source := ast.Source{
		Name:  "Schema",
		Input: *s,
	}
	docs, err := parser.ParseSchema(&source)
	if err.(*gqlerror.Error) != nil {
		log.Fatalln("Failed to parse given schema", err)
	}
	enumList := make(EnumList, 0)
	for _, def := range docs.Definitions {
		if def.Kind == "ENUM" {
			values := make([]string, len(def.EnumValues))
			for i, v := range def.EnumValues {
				values[i] = v.Name
			}
			enumList = append(enumList, &Enum{
				Name:   def.Name,
				Values: values,
			})
		}
	}
	objList := make(ObjectList, 0)
	connections := utils.PairList{}
	hasManyList := utils.PairList{}
	for _, def := range docs.Definitions {
		if def.Kind == "ENUM" {
			continue
		}
		fields := make(FieldList, len(def.Fields))
		for fi, field := range def.Fields {
			d := directive.NewDirective(field.Directives)
			if d != nil {
				if d.ConnectionType.IsManyToMany() {
					connection := &utils.Pair{
						First:  def.Name,
						Second: field.Type.Elem.NamedType,
					}
					if connections.HasSamePair(*connection) == false {
						connections = append(connections, connection)
					}
				} else if d.ConnectionType.IsHasMany() {
					hasMany := &utils.Pair{
						First:  field.Type.Elem.NamedType,
						Second: def.Name,
					}
					if hasManyList.HasSamePair(*hasMany) {
						log.Fatalln("@hasMany directive is defined but @manyToMany should be used in this case")
					} else {
						hasManyList = append(hasManyList, hasMany)
					}
				}
			}
			if field.Type.Elem == nil {
				fields[fi] = &Field{
					Name: field.Name,
					ReturnType: ReturnType{
						Name:       field.Type.NamedType,
						IsArray:    false,
						IsRequired: field.Type.NonNull,
						EnumList:   &enumList,
					},
					Directive: d,
				}
			} else {
				fields[fi] = &Field{
					Name: field.Name,
					ReturnType: ReturnType{
						Name:            field.Type.Elem.NamedType,
						IsArray:         true,
						IsRequired:      field.Type.NonNull,
						IsArrayRequired: field.Type.Elem.NonNull,
						EnumList:        &enumList,
					},
					Directive: d,
				}
			}
		}
		authDirective := def.Directives.ForName("auth")
		var authRules *directive.AuthRuleList
		if authDirective != nil {
			authRules = directive.NewAuthRuleList(*authDirective)
		}
		objList = append(objList, &Object{
			Name:      def.Name,
			Fields:    fields,
			AuthRules: authRules,
			Type:      NewObjectType(def.Name),
		})
	}
	return &Schema{
		Objects:     objList,
		Enums:       enumList,
		FilePath:    importPath,
		Connections: connections,
		HasManyList: hasManyList,
	}
}
