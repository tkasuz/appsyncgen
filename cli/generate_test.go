package cli

import (
	"testing"

	"github.com/kopkunka55/appsyncgen/codegen/api"
)

func TestGenerate(t *testing.T) {
	apiBuilder := api.NewAppSyncApiBuilder()
	apiBuilder.SetExportPath("./../resolvers")
	apiBuilder.SetSchema("./../testdata/schema.graphql")
	apiBuilder.AddDataSource("DYNAMODB", "appsync")
	apiBuilder.SetName("appsync")
	appSync := apiBuilder.Build()
	appSync.Export()
	appSync.Synth()
}
