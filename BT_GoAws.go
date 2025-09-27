/*
===============================================================================
AWS CDK INFRASTRUCTURE AS CODE - GO LEARNING PROJECT
===============================================================================

LEARNING PROGRESSION NOTES:
- Currently active: DynamoDB Table + Lambda Function (basic infrastructure)
- Commented out for later: API Gateway + IAM Permissions + API Routing
- To uncomment sections: Look for "TODO: Uncomment this section when learning..." comments

===============================================================================

Author: Ben Tran
Date: 22/09/2025
Description: This project demonstrates AWS CDK (Cloud Development Kit) using Go to define
             and deploy cloud infrastructure including Lambda functions, API Gateway,
             and DynamoDB tables. This is Infrastructure as Code (IaC) that allows
             you to manage AWS resources programmatically.

? TO RUN THIS APPLICATION LOCAL MACHINE:
```bash
# 1. Start DynamoDB Local (if using local DB)
   cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS"
   docker run -p 8000:8000 amazon/dynamodb-local

# 2. Run your Lambda function locally
cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS/lambda"
go run main.go
```
# 3. Create the table in Bash terminal
cd "/c/Users/Admin/Documents/GODEV/GO Courses/GO AWS/BT_GoAWS"
aws dynamodb create-table \
    --table-name BTtableGuestsInfo \
    --attribute-definitions AttributeName=username,AttributeType=S \
    --key-schema AttributeName=username,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000

# 4. AWS WORKAROUND >>>>> Add an item manually # Open Git Bash terminal
# Run the command to manually add an item (much easier in bash!)
aws dynamodb put-item \
    --table-name BTtableGuestsInfo \
    --item '{"username":{"S":"Phonenix !"},"password":{"S":"Password123 ?"}}' \
    --endpoint-url http://localhost:8000

# 5. View the table
aws dynamodb list-tables --endpoint-url http://localhost:8000
aws dynamodb scan --table-name BTtableGuestsInfo --endpoint-url http://localhost:8000


# ⚠️  Delete the table
aws dynamodb delete-table --table-name BTtableGuestsInfo --endpoint-url http://localhost:8000

///TODO - THIS IS FOR THE LAMBDA FUNCTION 2 DEPLOY TO AWS -FROM WINDOWS PLATFORM
# Run the fixed PowerShell script
.\Deploy_with_optimization.ps1

# NOTE: OPTIMIZATION -
$env:ENABLE_PERFORMANCE_MONITORING="true"
$env:OPTIMIZE_MEMORY="true"
$env:OPTIMIZE_DATABASE="true"
$env:ENABLE_CACHING="true"
$env:ENABLE_CLOUDWATCH="true"
$env:ENABLE_METRICS="true"
$env:ENABLE_LOGGING="true"
$env:ENABLE_HEALTH_CHECK="true"

# Step 1: Build Lambda Function (Linux)
cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS\lambda"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o bootstrap main.go
Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

# Step 2: Build CDK Application (Windows)
cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o BTgoAWS.exe BT_GoAws.go

# Step 3: Test Synthesis (FIXED - was missing)
cdk synth
cdk diff

# Step 4: Deploy New Stack
cdk deploy --yes --require-approval never

# Step 5: Test via AWS CLI (FIXED - better payload)
aws lambda invoke --function-name BTgoAWSstack-BTlambdaFunction8DF41906-ZlGN7I54bIaq --payload '{"path": "/register", "body": "{\"username\": \"testuser\", \"password\": \"password123\"}"}' response.json
Get-Content response.json

///NOTE - Test via API Gateway
cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS\lambda"
# Run all tests
go test -v ./...

# Run specific test categories
go test -v ./types
go test -v ./api
go test -v ./database

# Run integration tests
go test -v ./integration_test.go

# Run E2E tests
go test -v ./e2e_test.go
================================================
#DELETE TABLES
aws dynamodb list-tables
aws dynamodb delete-table --table-name userTable2

# Check Roll Back Stack Status
aws cloudformation describe-stacks --stack-name BTgoAWSstack --query 'Stacks[0].StackStatus' --output text
aws cloudformation describe-stack-events --stack-name BTgoAWSstack --query 'StackEvents[?ResourceStatus==`CREATE_FAILED` || ResourceStatus==`UPDATE_FAILED`].{Resource:LogicalResourceId,Status:ResourceStatus,Reason:ResourceStatusReason}' --output table

///TUB - work file
cd "C:\Users\Admin\Documents\GODEV\BTGo\BT_GoAWS"
go work init .
go work use ./lambda

LEARNING OBJECTIVES:
- AWS CDK (Cloud Development Kit) with Go
- Infrastructure as Code (IaC) concepts
- Lambda function deployment and configuration
- API Gateway setup with CORS and routing
- DynamoDB table creation and permissions
- CloudFormation stack management
- AWS resource integration and dependencies

/// REVIEW -
TOPICS COVERED:
- AWS CDK Go constructs and patterns
- Lambda function deployment with custom runtime
- API Gateway REST API with CORS configuration
- DynamoDB table with partition key setup
- IAM permissions and resource access
- CloudFormation stack synthesis
- Environment configuration and deployment

FEATURES:
- Serverless Lambda function with Go runtime
- RESTful API with multiple endpoints (/register, /login, /protected)
- DynamoDB table for user data storage
- CORS-enabled API Gateway for web applications
- Proper IAM permissions for Lambda-DynamoDB access
- Regional API Gateway endpoint configuration
- Request/response logging enabled

ARCHITECTURE:
- CDK Stack: Container for all AWS resources
- Lambda Function: Serverless compute for API logic
- API Gateway: HTTP API with CORS and routing
- DynamoDB: NoSQL database for user data
- IAM: Permissions for resource access

///LINK - https://us-east-1.console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks?filteringText=&filteringStatus=active&viewNested=true


///NOTE: $ cdk cli-telemetry --disable
$env:CDK_DISABLE_CLI_TELEMETRY="true"
go clean -cache
$ aws sts get-caller-identity
$ cdk deploy --require-approval never
$ aws cloudformation list-stacks --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE CREATE_FAILED UPDATE_FAILED
$ cdk deploy --verbose --debug
$ cdk synth
$ cdk deploy --require-approval never --concurrency 1
$ aws cloudformation describe-stacks --stack-name BTgoAWSstack
$ cdk acknowledge 34892
$ cdk ack 34892
$ cdk deploy --yes
$ aws cloudformation describe-stacks --stack-name BTgoAWSstack --query 'Stacks[0].StackStatus' --output text
$ cdk deploy --yes --no-rollback
$ go build -o bootstrap main.go
$ go build -o bootstrap main.go
$ aws dynamodb list-tables --query 'TableNames[?contains(@, `userTable`)]' --output table
$ aws cloudformation delete-stack --stack-name BTgoAWSstack
$ aws configure get region
$ cdk context
$ aws cloudformation execute-change-set --change-set-name cdk-deploy-change-set --stack-name BTgoAWSstack


*/

// Package main is the entry point for our AWS CDK application.
// This package defines the cloud infrastructure using Infrastructure as Code (IaC)
// principles, allowing us to manage AWS resources programmatically.
package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"               // Core CDK library for AWS resources and stacks
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway" // API Gateway constructs for REST APIs (commented out for later use)
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"   // DynamoDB constructs for NoSQL database
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"     // Lambda constructs for serverless functions
	"github.com/aws/constructs-go/constructs/v10"       // Core constructs library for building CDK apps
	"github.com/aws/jsii-runtime-go"                    // JavaScript interop runtime (needed for CDK)
)

const TABLE_NAME = "BTtableGuestsInfo" // AWS DynamoDB table name for user data
const LOCAL_TABLE_NAME = "LocoTable"   // Local testing table name
// BTgoAWSstackProps defines the properties for our CDK stack.
// It embeds awscdk.StackProps to inherit standard stack configuration options
// like environment, tags, and other AWS-specific settings.
type BTgoAWSstackProps struct {
	awscdk.StackProps // Embedded struct containing stack properties like environment, tags, etc.
}

// NewBTgoAWSstack creates a new CDK stack with the given scope, ID, and properties.
// This is where you define all your AWS resources (Lambda functions, API Gateway, DynamoDB, etc.)
// The stack acts as a container that groups related resources together for deployment.
//
// Parameters:
//   - scope: The parent construct (usually the CDK App)
//   - id: Unique identifier for this stack
//   - props: Optional stack properties (environment, tags, etc.)
//
// Returns:
//   - awscdk.Stack: The created stack containing all defined resources
func NewBTgoAWSstack(scope constructs.Construct, id string, props *BTgoAWSstackProps) awscdk.Stack {
	// Handle optional properties - if none provided, use default empty properties
	// This pattern allows for flexible stack configuration
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps // Extract the embedded StackProps
	}

	// Create the actual CDK stack - this is the container for all your AWS resources
	// The stack will be deployed as a CloudFormation stack in AWS
	stack := awscdk.NewStack(scope, &id, &sprops)

	// ========================================
	// DYNAMODB TABLE DEFINITION
	// ========================================
	// Create a DynamoDB table to store user data
	// DynamoDB is a NoSQL database that scales automatically and provides
	// single-digit millisecond latency at any scale
	table := awsdynamodb.NewTable(stack, jsii.String(TABLE_NAME), &awsdynamodb.TableProps{
		// Define the partition key (primary key) for the table
		// The partition key determines how data is distributed across partitions
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),          // Field name for the partition key
			Type: awsdynamodb.AttributeType_STRING, // Data type (STRING, NUMBER, BINARY)
		},
		TableName: jsii.String(TABLE_NAME), // Physical table name in AWS
	})

	// ========================================
	// LAMBDA FUNCTION DEFINITION
	// ========================================
	// Create an AWS Lambda function that will handle our API requests
	// Lambda is a serverless compute service that runs code without provisioning servers
	// NOTE: awslambda.functionprops has a lot of options, so you can customize the Lambda function to your needs
	myFunction := awslambda.NewFunction(stack, jsii.String("BTlambdaFunction"), &awslambda.FunctionProps{
		// Runtime environment for the Lambda function
		// PROVIDED_AL2023 uses Amazon Linux 2023 with Go runtime
		Runtime: awslambda.Runtime_PROVIDED_AL2023(), // Runtime environment for the Lambda function by FunctionProps

		// Source code for the Lambda function
		// FromAsset loads code from a ZIP file in the filesystem
		Code: awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil), // Source code for the Lambda function by FunctionProps
		// - zip file is created by running 'go build -o main main.go' and 'zip deployment.zip main' in the lambda directory

		// Entry point for the Lambda function
		// "main" refers to the main function in your Go code
		Handler: jsii.String("HandleAPIGatewayRequest"), // Entry point for the Lambda function by FunctionProps
		//FIXME - Your CDK is using Handler: jsii.String("main") which calls the main function,
		// but your main function is calling the wrong handler based on LocoOrLambda!

		// NOTE: Environment variables TABLE_NAME and LOCO_OR_LAMBDA are used to pass values from CDK to Lambda as Global
		// Environment variables for the Lambda function
		Environment: &map[string]*string{
			"TABLE_NAME":     jsii.String(TABLE_NAME), // AWS table: BTtableGuestsInfo get from default CDK constant TABLE_NAME
			"LOCO_OR_LAMBDA": jsii.String("99"),       // 0=Local, 1=Direct Lambda, 99=API Gateway
		},
	})

	// ========================================
	// API GATEWAY DEFINITION (COMMENTED OUT FOR LATER USE)
	// ========================================
	// TODO: Uncomment this section when learning API Gateway in class
	// Create an API Gateway REST API to expose our Lambda function as HTTP endpoints
	// API Gateway is a fully managed service that makes it easy to create, publish,
	// maintain, monitor, and secure APIs at any scale

	api := awsapigateway.NewRestApi(stack, jsii.String("myAPIGateway2"), &awsapigateway.RestApiProps{
		// CORS (Cross-Origin Resource Sharing) configuration
		// This allows web applications from different domains to access our API
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),           // Allowed request headers
			AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"), // Allowed HTTP methods
			AllowOrigins: jsii.Strings("*"),                                       // Allow requests from any origin (use specific domains in production)
		},

		// Deployment and logging configuration
		DeployOptions: &awsapigateway.StageOptions{
			//FIXME - You can remove that line or follow these steps to enable logging >> In the REVIEW section below
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO, // Enable INFO-level logging for debugging
		},

		// Endpoint configuration - determines where the API is accessible
		EndpointConfiguration: &awsapigateway.EndpointConfiguration{
			Types: &[]awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL}, // Regional endpoint for better performance
		},
	})

	// ========================================
	// API ROUTING AND INTEGRATION (COMMENTED OUT FOR LATER USE)
	// ========================================
	// TODO: Uncomment this section when learning API Gateway routing in class
	// Create a Lambda integration that connects API Gateway to our Lambda function
	// This integration handles the transformation between HTTP requests and Lambda events

	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

	// Define API routes and their HTTP methods
	// Each route corresponds to an endpoint in our Lambda function

	// /register endpoint - for user registration
	registerRoute := api.Root().AddResource(jsii.String("register"), nil)
	registerRoute.AddMethod(jsii.String("POST"), integration, nil) // POST method for creating new users
	//NOTE - ADDING NEW ENDPOINT - /register - for user registration
	//registerRoute := api.Root().AddResource(jsii.String("register/ID"), nil)
	// ID is the parameter for the endpoint - if needed

	// /login endpoint - for user authentication
	loginRoute := api.Root().AddResource(jsii.String("login"), nil)
	loginRoute.AddMethod(jsii.String("POST"), integration, nil) // POST method for user login

	// /protected endpoint - for authenticated access
	protectedRoute := api.Root().AddResource(jsii.String("protected"), nil)
	protectedRoute.AddMethod(jsii.String("GET"), integration, nil) // GET method for protected resources

	// ========================================
	// IAM PERMISSIONS (COMMENTED OUT FOR LATER USE)
	// ========================================
	// // Return the completed stack with all resources defined
	// // Note: Variables are used to avoid "declared and not used" errors
	// _ = table      // DynamoDB table (will be used when IAM permissions are uncommented)
	// _ = myFunction // Lambda function (will be used when API Gateway is uncommented)
	// TODO: Uncomment this section when learning IAM permissions in class
	// Grant the Lambda function read and write permissions to the DynamoDB table
	// This is essential for the Lambda function to access the database
	// CDK automatically creates the necessary IAM policies and roles

	// NOTE: GrantReadWriteData is a method that grants the Lambda function read and write permissions to the DynamoDB table
	// It is a method of the awsdynamodb.Table construct
	// It is used to grant the Lambda function read and write permissions to the DynamoDB table
	// It is used to grant the Lambda function read and write permissions to the DynamoDB table
	table.GrantReadWriteData(myFunction)

	return stack
}

// main() is the entry point of your CDK application.
// This function orchestrates the entire CDK deployment process.
func main() {
	// CRITICAL: Always close the JSII runtime when the program exits
	// JSII is the JavaScript interop layer that CDK uses - it needs proper cleanup
	defer jsii.Close()

	// Create a new CDK application - this is the root of your CDK app
	// Think of it as the "project" that contains all your stacks
	app := awscdk.NewApp(nil)

	// Create your stack and add it to the application
	// This is where you instantiate your stack with a name and configuration
	NewBTgoAWSstack(app, "BTgoAWSstack", &BTgoAWSstackProps{
		awscdk.StackProps{
			Env: env(), // Set the AWS environment (account/region) for deployment
		},
	})

	// SYNTHESIZE: Convert your CDK code into CloudFormation templates
	// This is the final step - it generates the actual AWS CloudFormation JSON/YAML
	// that will be deployed to AWS
	app.Synth(nil)
}

// env() determines the AWS environment (account+region) where your stack will be deployed.
// This is crucial for CDK to know which AWS account and region to target.
// For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// OPTION 1: Environment-agnostic (CURRENT SETTING)
	// If unspecified, this stack will be "environment-agnostic"
	// - Account/Region-dependent features and context lookups will not work
	// - Single synthesized template can be deployed anywhere
	// - Good for development and testing
	//---------------------------------------------------------------------------
	// return nil

	// OPTION 2: Specific environment (CURRENT SETTING)
	return &awscdk.Environment{
		Account: jsii.String("509505638450"), // Your AWS account ID
		Region:  jsii.String("us-east-1"),    // Your preferred AWS region
	}
}

//REVIEW - You can remove the DeployOptions or follow these instructions to enable logging:
/*
Go to the AWS IAM Console.
Navigate to Roles > Create role.
Select AWS Service, then choose API Gateway.
Choose the "API Gateway Push to CloudWatch Logs" policy.
Click Next, and name the role something like APIGatewayCloudWatchLogsRole.
Click Create role, and note down the ARN (Amazon Resource Name) of the role.
Then run:
aws apigateway update-account --patch-operations op=replace,path=/cloudwatchRoleArn,value=arn:aws:iam::id:role/APIGatewayCloudWatchLogsRole --profile profile
*/

/*REVIEW - COMPREHENSIVE CODE ANALYSIS
===============================================================================
EXPLANATION AND ANALYSIS OF AWS CDK INFRASTRUCTURE
===============================================================================

EXPLAINING THE CODE:
=========================
This AWS CDK application implements Infrastructure as Code (IaC) with the following architecture:

• Infrastructure Flow: CDK Code → CloudFormation Templates → AWS Resources
• Serverless Architecture: Lambda + API Gateway + DynamoDB
• CORS-enabled API: Web application ready with proper headers
• IAM Integration: Automatic permission management between services

STRUCTURE OF THE CODE:
===========================
1. Package Declaration & Imports:
   - AWS CDK core libraries (awscdk/v2)
   - Service-specific constructs (apigateway, dynamodb, lambda)
   - Core constructs and JSII runtime

2. Stack Definition:
   - BTgoAWSstackProps: Stack configuration structure
   - NewBTgoAWSstack: Main stack constructor with all resources

3. AWS Resources:
   - DynamoDB Table: User data storage with username partition key
   - Lambda Function: Go runtime with custom code deployment
   - API Gateway: REST API with CORS and regional endpoint
   - IAM Permissions: Automatic Lambda-DynamoDB access

4. API Endpoints:
   - /register: POST endpoint for user registration
   - /login: POST endpoint for user authentication
   - /protected: GET endpoint for authenticated access

WHY WE USE CDK:
===============
• Infrastructure as Code: Version control for infrastructure
• Type Safety: Compile-time checking of resource configurations
• Reusability: Share and reuse infrastructure components
• AWS Integration: Native integration with all AWS services
• CloudFormation: Leverages AWS's native deployment system

HOW IT WORKS:
=============
1. CDK synthesizes Go code into CloudFormation templates
2. CloudFormation deploys resources to AWS in correct order
3. API Gateway creates HTTP endpoints
4. Lambda function handles API requests
5. DynamoDB stores and retrieves user data
6. IAM ensures proper permissions between services

DEPLOYMENT PROCESS:
===================
1. Write CDK code defining infrastructure
2. Run 'cdk synth' to generate CloudFormation templates
3. Run 'cdk deploy' to create AWS resources
4. Test API endpoints with HTTP requests
5. Run 'cdk destroy' to clean up resources

IMPROVEMENT IDEAS:
==================
• Add environment-specific configurations (dev, staging, prod)
• Implement custom domain names for API Gateway
• Add CloudWatch alarms and monitoring
• Use AWS Secrets Manager for sensitive data
• Add VPC configuration for network isolation
• Implement API Gateway throttling and caching
• Add DynamoDB Global Secondary Indexes (GSI)
• Use AWS Certificate Manager for HTTPS
• Add CloudFront distribution for global CDN
• Implement AWS WAF for security
• Add X-Ray tracing for debugging
• Use AWS Systems Manager Parameter Store
• Add backup and disaster recovery
• Implement blue-green deployments

SECURITY CONSIDERATIONS:
========================
• IAM least privilege principle (already implemented)
• API Gateway authentication and authorization
• DynamoDB encryption at rest and in transit
• VPC endpoints for private API access
• AWS WAF for DDoS protection
• CloudTrail for audit logging
• Secrets management for sensitive data

COST OPTIMIZATION:
==================
• Use DynamoDB On-Demand billing for variable workloads
• Configure Lambda memory and timeout appropriately
• Use API Gateway caching to reduce Lambda invocations
• Monitor CloudWatch metrics for cost optimization
• Use AWS Free Tier effectively
• Implement auto-scaling policies

MONITORING AND OBSERVABILITY:
=============================
• CloudWatch Logs for Lambda function logs
• CloudWatch Metrics for performance monitoring
• X-Ray tracing for request flow analysis
• API Gateway execution logs
• DynamoDB CloudWatch metrics
• Custom CloudWatch dashboards

///TODO - ADDING FUNCTIONS RELATIONSHIPS -
arn:aws:iam::509505638450:role/APIGatewayCloudWatchLogsRole
aws apigateway update-account --patch-operations op=replace,path=/cloudwatchRoleArn,value=arn:aws:iam::509505638450:role/APIGatewayCloudWatchLogsRole --profile BTgoAWS_092025


CONCLUSION:
===========
This CDK application demonstrates a production-ready approach to defining
serverless infrastructure using Infrastructure as Code. The modular design,
proper resource configuration, and AWS best practices make it suitable for
real-world applications.

Key Takeaways:
- Infrastructure as Code provides version control and reproducibility
- CDK offers type safety and better developer experience than raw CloudFormation
- Proper IAM permissions are crucial for service integration
- CORS configuration enables web application integration
- Regional endpoints provide better performance than edge endpoints
- Always plan for monitoring, security, and cost optimization
- Use environment-specific configurations for different deployment stages

*/
