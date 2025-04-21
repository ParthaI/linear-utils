# ECR Enhanced Security Scanning

A comprehensive guide to automating Docker image security scanning with AWS ECR Enhanced Scanning.

## Overview

This guide provides a simple script and instructions for automating the process of:

1. Building a Docker image
2. Pushing it to Amazon ECR
3. Enabling enhanced security scanning
4. Retrieving vulnerability findings

## Prerequisites

- AWS CLI configured with appropriate permissions
- Docker installed and running
- An existing ECR repository or permissions to create one

## Quick Start

1. Download the [scan_ecr_enhanced.sh](scan_ecr_enhanced.sh) script
2. Make it executable: `chmod +x scan_ecr_enhanced.sh`
3. Run the script: `./scan_ecr_enhanced.sh`

## Script Explanation

The script performs the following operations:

### 1. Configuration Setup

```bash
REGION="ap-south-1"
REPO_NAME="test-enhanced-findings"
IMAGE_TAG="latest"
LOCAL_IMAGE_TAG="${REPO_NAME}:${IMAGE_TAG}"
```

The script defines configuration variables for the AWS region, repository name, and image tags.

### 2. ECR Authentication

```bash
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_IMAGE="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO_NAME}:${IMAGE_TAG}"
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"
```

The script retrieves your AWS account ID and authenticates with ECR.

### 3. Docker Image Creation

```bash
# Create Dockerfile if not present
if [ ! -f Dockerfile ]; then
  cat <<EOF > Dockerfile
FROM vulnerables/web-dvwa
EOF
fi

docker build -t $LOCAL_IMAGE_TAG .
docker tag $LOCAL_IMAGE_TAG $ECR_IMAGE
docker push $ECR_IMAGE
```

It creates a simple Dockerfile if one doesn't exist, builds the image, tags it, and pushes it to ECR.

### 4. Enhanced Scanning Configuration

```bash
aws ecr put-image-scanning-configuration \
  --repository-name $REPO_NAME \
  --image-scanning-configuration scanOnPush=true \
  --region $REGION
```

This enables enhanced scanning for the repository.

### 5. Scan Results Retrieval

```bash
# Wait for scan to complete
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
aws ecr describe-image-scan-findings \
  --repository-name $REPO_NAME \
  --image-id imageTag=$IMAGE_TAG \
  --region $REGION \
  --query 'imageScanFindings.findings[*].[name,severity,description]' \
  --output table
```

The script waits for the scan to complete and then retrieves and displays the scan findings.

## Best Practices

- Always use specific tags rather than `latest` in production
- Set up scan notifications to alert your team of critical vulnerabilities
- Use this script as part of your CI/CD pipeline
- Review and remediate findings promptly

## Troubleshooting

If you encounter issues:

1. Ensure your AWS CLI is properly configured
2. Check that your user has appropriate permissions for ECR operations
3. Verify that the ECR repository exists
4. Check Docker is running and can connect to ECR

## Additional Resources

- [AWS ECR Enhanced Scanning Documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
