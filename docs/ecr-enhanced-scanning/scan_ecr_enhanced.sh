#!/bin/bash

set -e

# Configurations
REGION="ap-south-1"
REPO_NAME="test-enhanced-findings"
IMAGE_TAG="latest"
LOCAL_IMAGE_TAG="${REPO_NAME}:${IMAGE_TAG}"

# Get AWS Account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Construct full ECR image URI
ECR_IMAGE="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO_NAME}:${IMAGE_TAG}"

echo "Logging into ECR..."
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

# Create Dockerfile if not present
if [ ! -f Dockerfile ]; then
cat <<EOF > Dockerfile
FROM vulnerables/web-dvwa
EOF
fi

echo "Building Docker image..."
docker build -t $LOCAL_IMAGE_TAG .

echo "Tagging image for ECR..."
docker tag $LOCAL_IMAGE_TAG $ECR_IMAGE

echo "Pushing image to ECR..."
docker push $ECR_IMAGE

echo "Enabling enhanced scanning for the repository..."
aws ecr put-image-scanning-configuration \
  --repository-name $REPO_NAME \
  --image-scanning-configuration scanOnPush=true \
  --region $REGION

# Wait for scan to complete
echo "Waiting for scan to complete..."
while true; do
  STATUS=$(aws ecr describe-image-scan-findings \
    --repository-name $REPO_NAME \
    --image-id imageTag=$IMAGE_TAG \
    --region $REGION \
    --query 'imageScanStatus.status' \
    --output text 2>/dev/null || echo "PENDING")

  echo "Current scan status: $STATUS"
  if [[ "$STATUS" == "COMPLETE" ]]; then
    break
  elif [[ "$STATUS" == "FAILED" ]]; then
    echo "Image scan failed."
    exit 1
  fi
  sleep 10
done

# Get scan findings
echo "Enhanced scan findings:"
aws ecr describe-image-scan-findings \
  --repository-name $REPO_NAME \
  --image-id imageTag=$IMAGE_TAG \
  --region $REGION \
  --query 'imageScanFindings.findings[*].[name,severity,description]' \
  --output table
