# Build Lambda function for Linux
build-lambda:
	@echo "Building Lambda function..."
	@cd lambda && GOOS=linux GOARCH=amd64 go build -o bootstrap
	@cd lambda && zip function.zip bootstrap
	@echo "Lambda function built successfully!"

# Build CDK application
build-cdk:
	@echo "Building CDK application..."
	@go build -o BTgoAWS.exe cdk5.go
	@echo "CDK application built successfully!"

# Deploy everything
deploy: build-lambda build-cdk
	@echo "Deploying to AWS..."
	@cdk deploy
	@echo "Deployment complete!"

# Destroy AWS resources (with confirmation)
destroy:
	@echo "Destroying AWS resources..."
	@cdk destroy
	@echo "AWS resources destroyed!"

# Force destroy AWS resources (no confirmation)
destroy-force:
	@echo "Force destroying AWS resources..."
	@cdk destroy --force
	@echo "AWS resources destroyed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f BTgoAWS.exe
	@rm -f lambda/bootstrap
	@rm -f lambda/function.zip
	@echo "Clean complete!"

# Destroy and clean everything
destroy-all: destroy clean
	@echo "Everything destroyed and cleaned!"

# Full destroy, clean, build and deploy
all: destroy clean deploy