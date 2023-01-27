package api

import (
	"appsyncgen/codegen/utils"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"log"
	"path/filepath"
)

type Props struct {
	awscdk.StackProps
	*AppSyncApi
}

var Code = `
export function request(ctx) {
    return {};
}
export function response(ctx) {
    return ctx.prev.result;
}
`

func Stack(scope constructs.Construct, id string, props *Props) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	graphqlApi := awsappsync.NewGraphqlApi(stack, jsii.String("AppSyncAPI"), &awsappsync.GraphqlApiProps{
		Name:        jsii.String(props.Name),
		XrayEnabled: jsii.Bool(false),
		Schema: awsappsync.NewSchemaFile(&awsappsync.SchemaProps{
			FilePath: jsii.String(props.Schema.FilePath),
		}),
	})
	table := awsdynamodb.NewTable(stack, jsii.String("DDBTable"), &awsdynamodb.TableProps{
		TableName:                  jsii.String(props.Name),
		BillingMode:                awsdynamodb.BillingMode_PAY_PER_REQUEST,
		ContributorInsightsEnabled: jsii.Bool(false),
		RemovalPolicy:              awscdk.RemovalPolicy_DESTROY,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})
	dataSource := awsappsync.NewDynamoDbDataSource(stack, jsii.String("AppSyncDDBDataSource"), &awsappsync.DynamoDbDataSourceProps{
		Api:                  graphqlApi,
		Name:                 jsii.String("AppSyncDynamoDBDataSource"),
		UseCallerCredentials: jsii.Bool(false),
		ReadOnlyAccess:       jsii.Bool(false),
		Table:                table,
	})
	table.GrantReadData(dataSource)
	table.GrantWriteData(dataSource)
	for _, r := range props.Resolvers {
		functions := make([]*string, len(r.Functions))
		for i, f := range r.Functions {
			code, err := utils.ReadFile(f)
			if err != nil {
				log.Fatalln("Failed to import code", err)
			}
			config := awsappsync.NewCfnFunctionConfiguration(stack, jsii.String(fmt.Sprintf("PipelineFun-%s-%s-%d", r.Object.Name, r.Name, i)), &awsappsync.CfnFunctionConfigurationProps{
				Name: jsii.String(fmt.Sprintf("%s_%s_%d", r.Object.Name, r.Name, i)),
				Code: code,
				Runtime: &map[string]interface{}{
					"name":           jsii.String("APPSYNC_JS"),
					"runtimeVersion": jsii.String("1.0.0"),
				},
				ApiId:          graphqlApi.ApiId(),
				DataSourceName: dataSource.Name(),
			})
			config.AddDependency(dataSource.Ds())
			functions[i] = config.AttrFunctionId()
		}
		resolver := awsappsync.NewCfnResolver(stack, jsii.String(fmt.Sprintf("Resolver-%s-%s", r.Object.Name, r.Name)), &awsappsync.CfnResolverProps{
			ApiId:     graphqlApi.ApiId(),
			Code:      &Code,
			TypeName:  jsii.String(r.Object.Name),
			FieldName: jsii.String(r.Name),
			Runtime: &map[string]interface{}{
				"name":           jsii.String("APPSYNC_JS"),
				"runtimeVersion": jsii.String("1.0.0"),
			},
			Kind: jsii.String("PIPELINE"),
			PipelineConfig: &map[string]interface{}{
				"functions": functions,
			},
		})
		resolver.AddDependency(dataSource.Ds())
	}
	return stack
}

func (a *AppSyncApi) Synth() {
	defer jsii.Close()
	app := awscdk.NewApp(&awscdk.AppProps{
		Outdir: jsii.String(filepath.Join(a.ExportPath, "cloudformation")),
	})
	Stack(app, "Stack", &Props{
		awscdk.StackProps{},
		a,
	})
	app.Synth(nil)
}
