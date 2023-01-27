package directive

import (
	"github.com/vektah/gqlparser/ast"
)

type Directive struct {
	AuthRules      *AuthRuleList
	ConnectionType *ConnectionType
}

func (d Directive) HasConnection() bool {
	if d.ConnectionType == nil {
		return false
	} else {
		return true
	}
}

func (d Directive) HasAuthRules() bool {
	if d.AuthRules == nil {
		return false
	} else {
		return true
	}
}

func NewDirective(l ast.DirectiveList) *Directive {
	if l == nil {
		return nil
	}
	dir := &Directive{}
	for _, it := range l {
		dir.ConnectionType = NewConnection(it.Name)
		dir.AuthRules = NewAuthRuleList(*it)
	}
	if dir.HasAuthRules() == false && dir.HasConnection() == false {
		return nil
	}
	return dir
}
