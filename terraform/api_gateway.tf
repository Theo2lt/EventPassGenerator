# api_gateway.tf

resource "aws_apigatewayv2_api" "my_serverless_app_api" {
  name          = "GoTicketAPI"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "my_serverless_app_lambda_integration" {
  api_id                 = aws_apigatewayv2_api.my_serverless_app_api.id
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
  integration_uri        = module.lambda_function.lambda_function_arn
}

resource "aws_apigatewayv2_route" "my_serverless_app_route" {
  api_id    = aws_apigatewayv2_api.my_serverless_app_api.id
  route_key = "POST /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.my_serverless_app_lambda_integration.id}"
}

resource "aws_apigatewayv2_stage" "my_serverless_app_stage" {
  api_id      = aws_apigatewayv2_api.my_serverless_app_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "my_serverless_app_api_gateway_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_function.lambda_function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.my_serverless_app_api.execution_arn}/*/*"
}

