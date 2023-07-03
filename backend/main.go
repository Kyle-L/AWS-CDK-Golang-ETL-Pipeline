package main

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	events "github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	glue "github.com/aws/aws-cdk-go/awscdkgluealpha/v2"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkWorkshopStackProps struct {
	awscdk.StackProps
	projectPrefix string
}

func NewCdkWorkshopStack(scope constructs.Construct, id string, props *CdkWorkshopStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// #################### 1. Date Migration ########################################
	// This portion of the stack creates the S3 bucket, DynamoDB table, and Glue job
	// that will be used to migrate the data from the S3 bucket to the DynamoDB table.
	// #############################################################################

	// Create a new S3 bucket
	bucket := s3.NewBucket(stack, jsii.String(`Bucket`), &s3.BucketProps{
		Encryption:         s3.BucketEncryption_S3_MANAGED,
		RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
		EventBridgeEnabled: jsii.Bool(true),
		AutoDeleteObjects:  jsii.Bool(true),
	})

	// Create a new Glue job.
	glueJob := glue.NewJob(stack, jsii.String("PythonETLJob"), &glue.JobProps{
		Executable: glue.JobExecutable_PythonEtl(&glue.PythonSparkJobExecutableProps{
			GlueVersion:   glue.GlueVersion_V3_0(),
			PythonVersion: glue.PythonVersion_THREE,
			Script:        glue.Code_FromAsset(jsii.String("glue/etl.py"), nil),
		}),
		Description: jsii.String("A simple Python ETL job"),
	})

	// Create a new DynamoDB table to store the results of the Glue job.
	table := dynamodb.NewTable(stack, jsii.String("Table"), &dynamodb.TableProps{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("id"),
			Type: dynamodb.AttributeType_NUMBER,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("accountNumber"),
			Type: dynamodb.AttributeType_STRING,
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		BillingMode:   dynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// Give glue access to read from the bucket.
	bucket.GrantRead(glueJob.GrantPrincipal(), nil)

	// Give glue job service role access to write to the table.
	table.GrantWriteData(glueJob.GrantPrincipal())

	// Create a new lambda function to trigger glue job.
	bundlingOptions := &awslambdago.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w"`)},
	}

	glueJobLambda := awslambdago.NewGoFunction(stack, jsii.String("GlueJobTriggerLambda"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("lambdas/glue-trigger"),
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
		Environment: &map[string]*string{
			"S3_KEY_PREFIX": jsii.String("input/"),
			"JOB_NAME":      glueJob.JobName(),
			"TABLE_NAME":    table.TableName(),
			"WORKERS":       jsii.String("8"),
		},
	})

	// Create EventBridge rule to trigger glue job lambda function.
	eventRule := events.NewRule(stack, jsii.String("S3ObjectCreated"), &events.RuleProps{
		EventPattern: &events.EventPattern{
			Source: &[]*string{jsii.String(`aws.s3`)},
			DetailType: &[]*string{
				jsii.String("Object Created"),
			},
			Detail: &map[string]interface{}{
				"bucket": map[string]interface{}{
					"name": &[]*string{bucket.BucketName()},
				},
			},
		},
	})
	eventRule.AddTarget(awseventstargets.NewLambdaFunction(glueJobLambda, nil))

	// Create IAM role for lambda to trigger glue job.
	glueJobLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"glue:StartJobRun",
		),
		Resources: jsii.Strings(
			*glueJob.JobArn(),
		),
	}))

	// Create Lambda to move finished files to archive folder.
	moveToArchiveLambda := awslambdago.NewGoFunction(stack, jsii.String("MoveToArchiveLambda"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("lambdas/move-to-archive"),
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
		Environment: &map[string]*string{
			"BUCKET_NAME": bucket.BucketName(),
		},
	})

	// Create EventBridge rule to trigger move to archive lambda function.
	eventRule = events.NewRule(stack, jsii.String("GlueStateChange"), &events.RuleProps{
		EventPattern: &events.EventPattern{
			Source: &[]*string{jsii.String(`aws.glue`)},
			DetailType: &[]*string{
				jsii.String("Glue Job State Change"),
			},
			Detail: &map[string]interface{}{
				"jobName": &[]*string{glueJob.JobName()},
				"state": &[]*string{
					jsii.String("SUCCEEDED"),
					jsii.String("FAILED"),
				},
			},
		},
	})
	eventRule.AddTarget(awseventstargets.NewLambdaFunction(moveToArchiveLambda, nil))

	// Create IAM role for lambda to move files to archive folder and give access to read from glue job.
	moveToArchiveLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"s3:ListBucket",
			"s3:GetObject",
			"s3:GetObjectTagging",
			"s3:PutObject",
			"s3:PutObjectTagging",
			"s3:DeleteObject",
			"glue:GetJobRun",
		),
		Resources: jsii.Strings(
			*bucket.BucketArn(),
			*bucket.BucketArn()+`/*`,
			*glueJob.JobArn(),
		),
	}))

	// #################### 2. API GATEWAY ########################################
	// This portion of the stack creates the API Gateway that will be used to
	// query the data in the DynamoDB table.
	// #############################################################################

	// Add global secondary index to the table for transactionDateTime and isFraud.
	// This is done so that we can query the table for all transactions that occurred
	// on a specific date and whether or not the transaction was fraudulent.
	table.AddGlobalSecondaryIndex(&dynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("isFraud-transactionDateTime-index"),
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("isFraud"),
			Type: dynamodb.AttributeType_STRING,
		},
		SortKey: &dynamodb.Attribute{
			Name: jsii.String("transactionDateTime"),
			Type: dynamodb.AttributeType_STRING,
		},
		ProjectionType: dynamodb.ProjectionType_ALL,
	})

	// Create a new lambda function to query the table.
	queryLambda := awslambdago.NewGoFunction(stack, jsii.String("QueryLambda"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("lambdas/dynamo-query"),
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
			"INDEX_NAME": jsii.String("isFraud-transactionDateTime-index"),
		},
	})

	// Grant the lambda function read access to the table.
	table.GrantReadData(queryLambda)

	// Create a new lambda function to update the table.
	updateLambda := awslambdago.NewGoFunction(stack, jsii.String("UpdateLambda"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("lambdas/dynamo-update"),
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
	})

	// Grant the lambda function read write access to the table.
	table.GrantReadWriteData(updateLambda)

	// Create a new API Gateway.
	api := apigateway.NewCfnApi(stack, jsii.String("API"), &apigateway.CfnApiProps{
		CorsConfiguration: &apigateway.CfnApi_CorsProperty{
			AllowHeaders:  jsii.Strings("*"),
			AllowMethods:  jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins:  jsii.Strings("*"),
			ExposeHeaders: jsii.Strings("*"),
		},
		Name:         jsii.String(strings.Join([]string{props.projectPrefix, `api`}, `-`)),
		ProtocolType: jsii.String("HTTP"),
	})

	// Get transactions route.
	queryIntegration := apigateway.NewCfnIntegration(stack, jsii.String("QueryIntegration"), &apigateway.CfnIntegrationProps{
		ApiId:                api.Ref(),
		IntegrationUri:       queryLambda.FunctionArn(),
		IntegrationType:      jsii.String("AWS_PROXY"),
		PayloadFormatVersion: jsii.String("1.0"),
	})
	apigateway.NewCfnRoute(stack, jsii.String("GetAllTransactionsResource"), &apigateway.CfnRouteProps{
		ApiId:             api.Ref(),
		AuthorizationType: jsii.String("NONE"),
		Target:            jsii.String("integrations/" + *queryIntegration.Ref()),
		RouteKey:          jsii.String("GET /transactions"),
	})
	queryLambda.AddPermission(jsii.String("QueryLambdaPermission"), &awslambda.Permission{
		Action:    jsii.String("lambda:InvokeFunction"),
		Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil),
		SourceArn: jsii.String("arn:aws:execute-api:" + *stack.Region() + ":" + *stack.Account() + ":" + *api.Ref() + "/*"),
	})

	// Update transaction route.
	updateIntegration := apigateway.NewCfnIntegration(stack, jsii.String("UpdateIntegration"), &apigateway.CfnIntegrationProps{
		ApiId:                api.Ref(),
		IntegrationUri:       updateLambda.FunctionArn(),
		IntegrationType:      jsii.String("AWS_PROXY"),
		PayloadFormatVersion: jsii.String("1.0"),
	})
	apigateway.NewCfnRoute(stack, jsii.String("UpdateTransactionsResource"), &apigateway.CfnRouteProps{
		ApiId:             api.Ref(),
		AuthorizationType: jsii.String("NONE"),
		Target:            jsii.String("integrations/" + *updateIntegration.Ref()),
		RouteKey:          jsii.String("PUT /transactions/{id}"),
	})
	updateLambda.AddPermission(jsii.String("QueryLambdaPermission"), &awslambda.Permission{
		Action:    jsii.String("lambda:InvokeFunction"),
		Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil),
		SourceArn: jsii.String("arn:aws:execute-api:" + *stack.Region() + ":" + *stack.Account() + ":" + *api.Ref() + "/*"),
	})

	// Create a new stage for the API Gateway.
	stage := apigateway.NewCfnStage(stack, jsii.String("Stage"), &apigateway.CfnStageProps{
		ApiId:      api.Ref(),
		StageName:  jsii.String("v1"),
		AutoDeploy: jsii.Bool(true),
	})

	// Output the API Gateway URL.
	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value: jsii.String("https://" + *api.Ref() + ".execute-api." + *stack.Region() + ".amazonaws.com/" + *stage.StageName()),
	})

	// #################### 3. FRONTEND ########################################
	// This portion of the stack creates the frontend that will be used to
	// query the data in the DynamoDB table.
	// #############################################################################

	// Create a new S3 bucket to store the frontend.
	bucketFrontend := s3.NewBucket(stack, jsii.String("FrontendBucket"), &s3.BucketProps{
		Encryption:           s3.BucketEncryption_S3_MANAGED,
		RemovalPolicy:        awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects:    jsii.Bool(true),
		WebsiteIndexDocument: jsii.String(`index.html`),

		// This is needed for single page apps.
		WebsiteErrorDocument: jsii.String(`index.html`),

		// Allow public access to the frontend.
		BlockPublicAccess: s3.NewBlockPublicAccess(&s3.BlockPublicAccessOptions{
			BlockPublicPolicy: jsii.Bool(false),
		}),
	})

	// A bucket policy is needed to allow public access to the frontend.
	// This is not recommended for production applications. Instead, you should
	// use CloudFront to serve the frontend.
	bucketFrontend.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"s3:GetObject",
		),
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewAnyPrincipal(),
		},
		Resources: jsii.Strings(
			*bucketFrontend.BucketArn()+`/*`,
		),
	}))

	// Output the bucket name.
	awscdk.NewCfnOutput(stack, jsii.String("BucketName"), &awscdk.CfnOutputProps{
		Value: bucketFrontend.BucketName(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	projectName := app.Node().TryGetContext(jsii.String("name")).(string)

	NewCdkWorkshopStack(app, projectName, &CdkWorkshopStackProps{
		awscdk.StackProps{
			// Here, we define the stack name as the name of the project.
			StackName: jsii.String(projectName),
		},
		projectName,
	})

	app.Synth(nil)
}
