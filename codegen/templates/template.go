package templates

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

type DynamoDBResolverTemplateData struct {
	PK          string
	SK          string
	Attributes  []string
	TableName   string
	Connections []DynamoDBResolverTemplateData
}

func NewTemplate() *template.Template {
	return template.New("")
}

func ImportTemplate(tmpl *template.Template, tmplPaths ...string) *template.Template {
	for _, p := range tmplPaths {
		tmpl = template.Must(tmpl.ParseGlob(p))
	}
	return tmpl
}

func ExecuteTemplate(data any, tmplPath string, file *os.File, tmpl *template.Template) {
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, tmplPath, data); err != nil {
		log.Fatalln(err)
	}
	result := b.String()
	if _, err := file.WriteString(fmt.Sprintln(result)); err != nil {
		log.Fatalln(err)
	}
}

func AddFunctionMap(tmpl *template.Template, funMap map[string]any) *template.Template {
	return tmpl.Funcs(funMap)
}
