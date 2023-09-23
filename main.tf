provider "aws" {
  region = "us-east-1"
  profile = "personal"
}

data "aws_region" "current" { }

resource "aws_lambda_function" "lifedashboard" {
  function_name    = "lifedashboard"
  filename         = "ldb.zip"
  handler          = "ldb"
  source_code_hash = sha256(filebase64("ldb.zip"))
  role             = aws_iam_role.lifedashboard.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 30
}

# empty roll
resource "aws_iam_role" "lifedashboard" {
  name               = "lifedashboard"
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }
}
POLICY
}

# Allow API gateway to invoke the lifedashboard Lambda function.
resource "aws_lambda_permission" "lifedashboard" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lifedashboard.arn
  principal     = "apigateway.amazonaws.com"
}

resource "aws_api_gateway_resource" "lifedashboard" {
  rest_api_id = aws_api_gateway_rest_api.lifedashboard.id
  parent_id   = aws_api_gateway_rest_api.lifedashboard.root_resource_id
  path_part   = "lifedashboard"
}

resource "aws_api_gateway_rest_api" "lifedashboard" {
  name = "lifedashboard"
  binary_media_types = [ "*/*" ]
}

#           GET
# Internet -----> API Gateway
resource "aws_api_gateway_method" "lifedashboard" {
  rest_api_id   = aws_api_gateway_rest_api.lifedashboard.id
  resource_id   = aws_api_gateway_resource.lifedashboard.id
  http_method   = "GET"
  authorization = "NONE"
}
resource "aws_api_gateway_method_settings" "lifedashboard" {
  rest_api_id = aws_api_gateway_rest_api.lifedashboard.id
  stage_name  = aws_api_gateway_stage.v1.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }
}

#              POST
# API Gateway ------> Lambda
# For Lambda the method is always POST and the type is always AWS_PROXY.
#
# The date 2015-03-31 in the URI is just the version of AWS Lambda.
resource "aws_api_gateway_integration" "lifedashboard" {
  rest_api_id             = aws_api_gateway_rest_api.lifedashboard.id
  resource_id             = aws_api_gateway_resource.lifedashboard.id
  http_method             = aws_api_gateway_method.lifedashboard.http_method
  type = "AWS_PROXY"
  integration_http_method = "POST"
  passthrough_behavior    = "WHEN_NO_MATCH"
  uri                     = aws_lambda_function.lifedashboard.invoke_arn
  content_handling = "CONVERT_TO_BINARY"
}

resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = aws_api_gateway_rest_api.lifedashboard.id
  resource_id = aws_api_gateway_resource.lifedashboard.id
  http_method = aws_api_gateway_method.lifedashboard.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration_response" "lifedashboard" {
  rest_api_id      = aws_api_gateway_rest_api.lifedashboard.id
  resource_id      = aws_api_gateway_resource.lifedashboard.id
  http_method      = aws_api_gateway_method.lifedashboard.http_method
  status_code      = aws_api_gateway_method_response.response_200.status_code
  content_handling = "CONVERT_TO_BINARY"

  depends_on = [
    aws_api_gateway_integration.lifedashboard
  ]
}

# define the URL of the API Gateway.
resource "aws_api_gateway_deployment" "lifedashboard_v1" {
  depends_on = [
    aws_api_gateway_integration.lifedashboard
  ]
  rest_api_id = aws_api_gateway_rest_api.lifedashboard.id
}

# Set the generated URL as an output. Run `terraform output url` to get this.
output "url" {
  value = "${aws_api_gateway_deployment.lifedashboard_v1.invoke_url}${aws_api_gateway_resource.lifedashboard.path}"
}

resource "aws_api_gateway_stage" "v1" {
  stage_name    = "v1"
  rest_api_id   = aws_api_gateway_rest_api.lifedashboard.id
  deployment_id = aws_api_gateway_deployment.lifedashboard_v1.id
  xray_tracing_enabled = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.lifedashboard_api_gateway.arn
    format          = "$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] \"$context.httpMethod $context.resourcePath $context.protocol\" $context.status $context.responseLength $context.requestId"
  }
  depends_on = [aws_cloudwatch_log_group.lifedashboard_api_gateway]
}


resource "aws_cloudwatch_log_group" "lifedashboard_api_gateway" {
  name = "lifedashboard_api_gateway"
}









# Allow AG to log to cloudwatch
resource "aws_api_gateway_account" "demo" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch.arn
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "cloudwatch" {
  name               = "api_gateway_cloudwatch_global"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "cloudwatch" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
    ]

    resources = ["*"]
  }
}
resource "aws_iam_role_policy" "cloudwatch" {
  name   = "default"
  role   = aws_iam_role.cloudwatch.id
  policy = data.aws_iam_policy_document.cloudwatch.json
}
