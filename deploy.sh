#!/bin/bash

# Veritas Deployment Script
# This script gets the latest commit SHA and updates K8s YAML files and deploys them

set -e  # Stop on any error

echo "🚀 Starting Veritas Deployment Script..."

# 1. Get the latest commit SHA
NEW_SHA=$(git rev-parse HEAD)
echo "📝 Latest commit SHA: $NEW_SHA"

# 2. ACR Login Server
ACR_SERVER="veritasacr.azurecr.io"

# 3. Service list
SERVICES=("creator-service" "redirector-service" "analytics-service" "frontend-service")

echo "🔄 Updating image tags..."

# 4. Update image tag in YAML file for each service
for service in "${SERVICES[@]}"; do
    yaml_file="k8s/${service}.yaml"
    
    if [ -f "$yaml_file" ]; then
        # Map service name to image name (frontend-service uses 'frontend' image)
        image_name=$service
        if [ "$service" = "frontend-service" ]; then
            image_name="frontend"
        fi
        
        # Replace old image tag with new SHA
        sed -i "s|image: ${ACR_SERVER}/veritas/${image_name}:.*|image: ${ACR_SERVER}/veritas/${image_name}:${NEW_SHA}|g" "$yaml_file"
        echo "✅ ${service} YAML file updated"
    else
        echo "⚠️  ${yaml_file} file not found, skipping..."
    fi
done

echo "🔄 Deploying to Kubernetes..."

# 5. Apply to Kubernetes
for service in "${SERVICES[@]}"; do
    yaml_file="k8s/${service}.yaml"
    
    if [ -f "$yaml_file" ]; then
        echo "📦 Deploying ${service}..."
        kubectl apply -f "$yaml_file"
    fi
done

echo "🎉 Deployment completed!"
echo "📊 Check pod status:"
echo "   kubectl get pods"
echo "   kubectl logs -l app=creator-service" 