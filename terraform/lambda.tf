# lambda.tf

# https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.17.0"

  function_name = "go-ticket-api"
  description   = "Lambda function to execute go-ticket-api"
  image_uri     = "${aws_ecr_repository.app.repository_url}:latest"
  architectures = ["arm64"]
  memory_size   = 1024
  timeout       = 200

  package_type   = "Image"
  create_package = false

  lambda_role = aws_iam_role.lambda_execution_role.arn
}