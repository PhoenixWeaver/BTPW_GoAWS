Options:
0: Local Testing
LocoOrLambda := 0
Local testing with direct JSON
Uses HandleRequest
1: Simple HandleRequest
LocoOrLambda := 1
AWS Lambda with direct JSON
Uses HandleRequest
Default: API Gateway (Production)
LocoOrLambda := 99 (or any other value)
AWS Lambda with API Gateway
Uses RegisterUser with full app structure
Ready for middleware implementation


>>>> TEST option 1:
Expected Result:
✅ Lambda returns: "Successfully processed user: Phoenix ! with password: Password123 ?"
✅ DynamoDB shows: New item with username "Phoenix" and hashed password
✅ No errors in CloudWatch logs