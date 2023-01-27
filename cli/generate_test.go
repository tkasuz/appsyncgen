package cli

import (
	"appsyncgen/codegen/api"
	"testing"
)

func TestGenerate(t *testing.T) {
	apiBuilder := api.NewAppSyncApiBuilder()
	apiBuilder.SetExportPath("./../resolvers")
	apiBuilder.SetTemplates("./../codegen/templates")
	apiBuilder.SetSchema("./../testdata/schema.graphql")
	apiBuilder.AddDataSource("DYNAMODB", "appsync")
	apiBuilder.SetName("appsync")
	appSync := apiBuilder.Build()
	appSync.Synth()
	appSync.Export()
}
