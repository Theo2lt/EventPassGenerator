# output.tf

output "api_endpoint" {
  description = "The endpoint of the API Gateway"
  value       = aws_apigatewayv2_api.my_serverless_app_api.api_endpoint
}
