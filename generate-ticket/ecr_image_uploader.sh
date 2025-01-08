#!/bin/bash

set -euo pipefail

# Variables
AWS_REGION="eu-west-1"
ECR_REPO_URL="533721337392.dkr.ecr.eu-west-1.amazonaws.com/go-ticket-api"
IMAGE_NAME="lambda-golang"
TAG=${1:-latest}  # Utiliser le premier argument comme tag, ou "latest" par défaut
PLATFORM="linux/arm64"  # Spécifiez la plateforme cible : linux/amd64 ou linux/arm64
LAMBDA_FUNCTION_NAME="go-ticket-api"

# Fonction pour afficher un message d'erreur et quitter
error_exit() {
    echo "Erreur : $1" >&2
    exit 1
}

# Connexion à Amazon ECR
echo "Connexion à ECR..."
aws ecr get-login-password --region "$AWS_REGION" | \
    docker login --username AWS --password-stdin "$ECR_REPO_URL" || \
    error_exit "Échec de la connexion à ECR."

# Construction de l'image Docker
echo "Construction de l'image Docker avec le tag '$TAG' pour la plateforme '$PLATFORM'..."
docker build --provenance=false --platform "$PLATFORM" -t "$IMAGE_NAME:$TAG" . || \
    error_exit "Échec de la construction de l'image Docker."

# Tag de l'image Docker
echo "Application du tag à l'image Docker..."
docker tag "$IMAGE_NAME:$TAG" "$ECR_REPO_URL:$TAG" || \
    error_exit "Échec de l'application du tag à l'image Docker."

# Envoi de l'image Docker à ECR
echo "Envoi de l'image Docker à ECR..."
docker push "$ECR_REPO_URL:$TAG" || \
    error_exit "Échec de l'envoi de l'image Docker à ECR."

# Mise à jour de la fonction Lambda en mode silencieux
echo "Mise à jour de la fonction AWS Lambda..."
if ! aws lambda update-function-code \
    --function-name "$LAMBDA_FUNCTION_NAME" \
    --image-uri "$ECR_REPO_URL:$TAG" > /dev/null 2>&1; then
    error_exit "Échec de la mise à jour de la fonction Lambda."
fi

echo "L'image Docker '$IMAGE_NAME:$TAG' a été poussée avec succès et la fonction Lambda mise à jour !"
