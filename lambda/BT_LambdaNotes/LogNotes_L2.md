The issue is clear now. Your AWS user PhoenixWeaver doesn't have the necessary permissions for CDK operations. Here's what you need to do:
Solution: Add Required IAM Permissions
You need to add the following permissions to your PhoenixWeaver user in the AWS IAM console:
Required Permissions for CDK:
CloudFormation permissions:


//FIXME - TROUBLESHOOTING - 
Quick Solution: Use AdministratorAccess (for learning/testing)
For learning purposes, you can temporarily attach the AdministratorAccess policy to your user, which includes all the permissions needed for CDK operations.
Steps to Fix:
Go to AWS IAM Console:
Navigate to https://console.aws.amazon.com/iam/
Go to "Users" → "PhoenixWeaver"
Add Permissions:
Click "Add permissions" → "Attach policies directly"
Search for and select AdministratorAccess (for learning)
OR create a custom policy with the specific permissions listed above
After adding permissions, try again:


Solution 2: The real issue is the AWS permissions
The error you're seeing is actually related to the AWS permissions issue we discussed earlier. The CDK is trying to run but failing due to insufficient permissions.
Let me update the todo and provide you with the correct solution:
The issue is not with the executable location. The cdk.json configuration is correct - it should run go run cdk2.go from the root directory. The real problem is the AWS permissions issue we identified earlier.
Here's what you need to do:
1. Fix AWS Permissions (This is the main issue)
Your AWS user PhoenixWeaver needs additional permissions. Go to the AWS IAM console and add these permissions:
Option A: Quick Fix (for learning/testing):
Attach the AdministratorAccess policy to your PhoenixWeaver user
Option B: Specific Permissions:
Add these permissions to your user:
cloudformation:*
ssm:*
s3:*
iam:*
2. After fixing permissions, try these commands:
3. The lambda folder structure is correct:
Your lambda folder contains:
bootstrap - The Lambda executable (this is correct)
main - Another executable (this is correct)
function.zip - The deployment package (this is what CDK uses)
The CDK code correctly references lambda/function.zip in your cdk2.go file.
The main issue is AWS permissions, not the executable location. Once you fix the permissions in the AWS IAM console, your CDK deployment should work properly.

//TODO - 
Perfect! I can see from the IAM console that your `PhoenixWeaver` user already has `AdministratorAccess` through the `AWS_Phoenix_CloudCrews` group. This means you have all the necessary permissions for CDK operations.

Since you have the proper permissions, let's try to bootstrap and deploy your CDK stack. The issue might be that the CDK bootstrap stack doesn't exist yet in your AWS account.

Let me help you run the commands to bootstrap and deploy:

## **Step 1: Bootstrap CDK Environment**

First, let's bootstrap the CDK environment. This creates the necessary AWS resources that CDK needs to deploy your stack:

```powershell
# Navigate to your CDK project root
cd "C:\Users\Admin\Documents\GODEV\GO Courses\GO AWS\CDK2"

# Bootstrap CDK (this sets up the necessary AWS resources)
cdk bootstrap
```

//SECTION - 
I see the issue now! The error shows that the Go executable is not being found in the PATH. This is a Windows-specific issue where the temporary executable created by `go run` is not accessible. Let me help you fix this.

The problem is that `go run` creates a temporary executable in a temp directory, but Windows can't find it in the PATH. Let's solve this by building the executable first and then running it.

## **Solution: Build the executable first**

Instead of using `go run cdk2.go`, let's build the executable and then run it:

### **Step 1: Build the CDK application**

```powershell
# Navigate to your CDK project root
cd "C:\Users\Admin\Documents\GODEV\GO Courses\GO AWS\CDK2"

# Build the executable
go build -o cdk2.exe cdk2.go
```

### **Step 2: Update cdk.json to use the built executable**

We need to update your `cdk.json` file to use the built executable instead of `go run`:

```json:cdk.json
{
 "app": "./cdk2.exe",
  "watch": {
    "include": [
      "**"
    ],
    "exclude": [
      "README.md",
      "cdk*.json",
      "go.mod",
      "go.sum",
      "**/*test.go"
    ]
  },
  // ... rest of the context remains the same
}
```

### **Step 3: Now try CDK commands**

After building and updating the configuration:

```powershell
# Bootstrap CDK
cdk bootstrap

# Deploy your infrastructure
cdk deploy
```

## **Alternative Solution: Use PowerShell execution policy**

If you prefer to keep using `go run`, you might need to adjust PowerShell execution policy:

```powershell
# Check current execution policy
Get-ExecutionPolicy

# If needed, set execution policy (run as Administrator)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## **Why this happens:**

- On Windows, `go run` creates temporary executables in `%TEMP%` directories
- These temporary paths are not in the system PATH
- The CDK process can't find the executable to run
- Building the executable first puts it in your project directory where it's accessible

Try the build approach first - it's the most reliable solution for Windows CDK deployments with Go.

//FIXME - UNINSTALL THE AWSCLI by PYTHON > INSTALL THE PRIVATE CONSOLE VERSION >>> Restart > let WSL.exe run > TURN OFF AVAST ANTIVIRUS

FIXED :
Stack Take2GoCdkStack
IAM Statement Changes
┌───┬──────────────────────────────────────┬────────┬─────────────────────────────┬──────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│   │ Resource                             │ Effect │ Action                      │ Principal
           │ Condition
                                                                    │
├───┼──────────────────────────────────────┼────────┼─────────────────────────────┼──────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ + │ ${myLambdaFunction2.Arn}             │ Allow  │ lambda:InvokeFunction       │ Service:apigateway.amazonaws.com


Do you wish to deploy these changes (y/n)?
