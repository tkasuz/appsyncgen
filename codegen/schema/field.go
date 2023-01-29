package schema

import (
	"strings"

	"github.com/kopkunka55/appsyncgen/codegen/directive"
)

type ReturnType struct {
	*EnumList
	Name            string `json:"name"`
	IsArray         bool   `json:"isArray"`
	IsRequired      bool   `json:"-"`
	IsArrayRequired bool   `json:"-"`
}

func (r ReturnType) IsPrimitive() bool {
	if r.Name == "ID" || r.Name == "String" || r.Name == "Int" || r.Name == "Float" || r.Name == "Boolean" {
		return true
	} else if r.Name == "AWSDate" || r.Name == "AWSTime" || r.Name == "AWSDateTime" || r.Name == "AWSTimestamp" ||
		r.Name == "AWSEmail" || r.Name == "AWSJSON" || r.Name == "AWSPhone" || r.Name == "AWSURL" || r.Name == "AWSIPAddress" {
		return true
	} else {
		return r.EnumList.IsEnum(r.Name)
	}
}

type Field struct {
	Name       string
	ReturnType ReturnType
	Directive  *directive.Directive
}

func (f Field) IsConnection() (*string, bool) {
	if index := strings.Index(f.Name, "ID"); index != -1 {
		name := f.Name[:index]
		return &name, true
	} else {
		return nil, false
	}
}

type FieldList []*Field

func (l FieldList) ForName(name string) *Field {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

func (l FieldList) Names() []string {
	names := make([]string, len(l))
	for i, f := range l {
		names[i] = f.Name
	}
	return names
}

func (l FieldList) HasConnection() bool {
	for _, f := range l {
		if _, ok := f.IsConnection(); ok {
			return true
		}
	}
	return false
}
