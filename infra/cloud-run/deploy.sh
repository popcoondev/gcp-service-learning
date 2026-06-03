#!/usr/bin/env bash

set -euo pipefail

PROJECT_ID="${PROJECT_ID:-gcp-service-learning}"
REGION="${REGION:-asia-northeast1}"
REPOSITORY="${REPOSITORY:-gcp-service-learning}"
SERVICE_NAME="${SERVICE_NAME:-order-api}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
RUNTIME_SERVICE_ACCOUNT="${RUNTIME_SERVICE_ACCOUNT:-order-api-runtime@${PROJECT_ID}.iam.gserviceaccount.com}"
IMAGE_URI="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${SERVICE_NAME}:${IMAGE_TAG}"

echo "Project ID: ${PROJECT_ID}"
echo "Region: ${REGION}"
echo "Repository: ${REPOSITORY}"
echo "Service Name: ${SERVICE_NAME}"
echo "Image URI: ${IMAGE_URI}"
echo "Runtime Service Account: ${RUNTIME_SERVICE_ACCOUNT}"

echo "==> Enabling required APIs"
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  iam.googleapis.com \
  logging.googleapis.com \
  monitoring.googleapis.com \
  --project "${PROJECT_ID}"

echo "==> Configuring Docker auth"
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

echo "==> Building image"
gcloud builds submit \
  --project "${PROJECT_ID}" \
  --tag "${IMAGE_URI}"

echo "==> Deploying Cloud Run service"
gcloud run deploy "${SERVICE_NAME}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image "${IMAGE_URI}" \
  --service-account "${RUNTIME_SERVICE_ACCOUNT}" \
  --allow-unauthenticated \
  --port 8080 \
  --timeout 300 \
  --min-instances 0 \
  --max-instances 10

echo "==> Deployment finished"
gcloud run services describe "${SERVICE_NAME}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --format="value(status.url)"
