# ecr.tf

resource "aws_ecr_repository" "app" {
  name                 = "go-ticket-api"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "null_resource" "docker_build_and_push" {
  depends_on = [aws_ecr_repository.app]

  provisioner "local-exec" {
    command = <<EOT
      # Authenticate Docker with ECR
      aws ecr get-login-password --region ${var.region} | docker login --username AWS --password-stdin ${aws_ecr_repository.app.repository_url}

      # Build the Docker image
      docker build -t ${var.docker_image_name}:${var.docker_image_tag} ../generate-ticket --build-arg REGION=${var.region}

      # Tag the Docker image for ECR
      docker tag ${var.docker_image_name}:${var.docker_image_tag} ${aws_ecr_repository.app.repository_url}:${var.docker_image_tag}

      # Push the Docker image to ECR
      docker push ${aws_ecr_repository.app.repository_url}:${var.docker_image_tag}
    EOT
  }
}
