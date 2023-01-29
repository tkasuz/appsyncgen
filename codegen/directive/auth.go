package directive

import (
	"fmt"
	"log"

	"github.com/vektah/gqlparser/ast"
)

type AuthProvider string

const (
	APIKEY  AuthProvider = "apiKey"
	IAM     AuthProvider = "iam"
	OIDC    AuthProvider = "oidc"
	COGNITO AuthProvider = "userPools"
	LAMBDA  AuthProvider = "lambda"
)

func (a AuthProvider) IsOidc() bool {
	return a == OIDC
}

func (a AuthProvider) IsIam() bool {
	return a == IAM
}

func (a AuthProvider) IsApiKey() bool {
	return a == APIKEY
}

func (a AuthProvider) IsCognito() bool {
	return a == COGNITO
}

func (a AuthProvider) IsLambda() bool {
	return a == LAMBDA
}

func IsAuthProvider(v string) bool {
	return v == string(APIKEY) || v == string(IAM) || v == string(OIDC) || v == string(COGNITO) || v == string(LAMBDA)
}

type AuthProviderList []*AuthProvider

func (l AuthProviderList) Intersection(ll AuthProviderList) (lll AuthProviderList) {
	m := make(map[AuthProvider]bool)
	for _, it := range l {
		m[*it] = true
	}
	for _, it := range ll {
		if _, ok := m[*it]; ok {
			lll = append(lll, it)
		}
	}
	return
}

func (l AuthProviderList) Union(ll AuthProviderList) (lll AuthProviderList) {
	m := make(map[AuthProvider]bool)
	for _, it := range l {
		m[*it] = true
		lll = append(lll, it)
	}
	for _, it := range ll {
		if _, ok := m[*it]; !ok {
			lll = append(lll, it)
		}
	}
	return
}

type ModelOperation string

const (
	CREATE ModelOperation = "create"
	UPDATE ModelOperation = "update"
	DELETE ModelOperation = "delete"
	READ   ModelOperation = "read"
)

func IsModelOperation(v string) bool {
	return v == string(CREATE) || v == string(UPDATE) || v == string(DELETE) || v == string(READ)
}

type AuthRule struct {
	AuthProvider    AuthProvider
	ModelOperations []ModelOperation
}

func (a AuthRule) HasCreate() bool {
	for _, it := range a.ModelOperations {
		if it == CREATE {
			return true
		}
	}
	return false
}

func (a AuthRule) HasUpdate() bool {
	for _, it := range a.ModelOperations {
		if it == UPDATE {
			return true
		}
	}
	return false
}

func (a AuthRule) HasRead() bool {
	for _, it := range a.ModelOperations {
		if it == READ {
			return true
		}
	}
	return false
}

func (a AuthRule) HasDelete() bool {
	for _, it := range a.ModelOperations {
		if it == DELETE {
			return true
		}
	}
	return false
}

type AuthRuleList []*AuthRule

func NewAuthRuleList(d ast.Directive) *AuthRuleList {
	if d.Name == "auth" {
		authRuleBuilder := NewAuthRuleBuilder()
		rules := d.Arguments.ForName("rules").Value.Children
		authRules := make(AuthRuleList, len(rules))
		for i, rule := range rules {
			provider := rule.Value.Children.ForName("provider")
			if provider == nil {
				log.Fatalln("provider should be given in @auth directive")
			}
			err := authRuleBuilder.SetAuthProvider(provider.Raw)
			if err != nil {
				log.Fatalln(err)
			}
			operations := rule.Value.Children.ForName("operations")
			if operations != nil {
				newOperations := make([]string, len(operations.Children))
				for oi, op := range operations.Children {
					newOperations[oi] = op.Value.Raw
				}
				err := authRuleBuilder.SetModelOperations(newOperations)
				if err != nil {
					log.Fatalln(err)
				}
			}
			authRule := authRuleBuilder.Build()
			authRules[i] = authRule
		}
		return &authRules
	}
	return nil
}

func (l AuthRuleList) ForProvider(v string) *AuthRule {
	if IsAuthProvider(v) {
		for _, it := range l {
			if it.AuthProvider == AuthProvider(v) {
				return it
			}
		}
	}
	return nil
}

func (l AuthRuleList) HasApiKey() bool {
	for _, it := range l {
		if it.AuthProvider == APIKEY {
			return true
		}
	}
	return false
}

func (l AuthRuleList) HasOidc() bool {
	for _, it := range l {
		if it.AuthProvider == OIDC {
			return true
		}
	}
	return false
}

func (l AuthRuleList) HasCognito() bool {
	for _, it := range l {
		if it.AuthProvider == COGNITO {
			return true
		}
	}
	return false
}

func (l AuthRuleList) HasIAM() bool {
	for _, it := range l {
		if it.AuthProvider == IAM {
			return true
		}
	}
	return false
}

func (l AuthRuleList) HasLambda() bool {
	for _, it := range l {
		if it.AuthProvider == LAMBDA {
			return true
		}
	}
	return false
}

func (l AuthRuleList) ProvidersAllowedToCreate() AuthProviderList {
	providers := make([]*AuthProvider, 0)
	for _, it := range l {
		if it.HasCreate() {
			providers = append(providers, &it.AuthProvider)
		}
	}
	return providers
}

func (l AuthRuleList) ProvidersAllowedToRead() AuthProviderList {
	providers := make([]*AuthProvider, 0)
	for _, it := range l {
		if it.HasRead() {
			providers = append(providers, &it.AuthProvider)
		}
	}
	return providers
}

func (l AuthRuleList) ProvidersAllowedToUpdate() AuthProviderList {
	providers := make([]*AuthProvider, 0)
	for _, it := range l {
		if it.HasUpdate() {
			providers = append(providers, &it.AuthProvider)
		}
	}
	return providers
}

func (l AuthRuleList) ProvidersAllowedToDelete() AuthProviderList {
	providers := make([]*AuthProvider, 0)
	for _, it := range l {
		if it.HasDelete() {
			providers = append(providers, &it.AuthProvider)
		}
	}
	return providers
}

type AuthRuleBuilder struct {
	AuthProvider    AuthProvider
	ModelOperations []ModelOperation
}

func NewAuthRuleBuilder() *AuthRuleBuilder {
	return &AuthRuleBuilder{
		ModelOperations: []ModelOperation{
			CREATE, READ, UPDATE, DELETE,
		},
	}
}

func (a *AuthRuleBuilder) SetAuthProvider(authProvider string) error {
	if IsAuthProvider(authProvider) {
		a.AuthProvider = AuthProvider(authProvider)
	} else {
		return fmt.Errorf("auth provider should be one of [apiKey, iam, oidc, userPools], but %s was given", authProvider)
	}
	return nil
}

func (a *AuthRuleBuilder) SetModelOperations(operations []string) error {
	modelOperations := make([]ModelOperation, len(operations))
	for i, op := range operations {
		if IsModelOperation(op) {
			modelOperations[i] = ModelOperation(op)
		} else {
			return fmt.Errorf("model operation should be one of [create, read, update, delete], but %s was given", op)
		}
	}
	a.ModelOperations = modelOperations
	return nil
}

func (a *AuthRuleBuilder) Build() *AuthRule {
	return &AuthRule{
		AuthProvider:    a.AuthProvider,
		ModelOperations: a.ModelOperations,
	}
}
