#!/bin/bash

# Veritas Deployment Script
# Bu script en son commit SHA'yı alıp K8s YAML dosyalarını günceller ve deploy eder

set -e  # Stop on any error

echo "🚀 Veritas Deployment Script Başlıyor..."

# 1. En son commit SHA'yı al
NEW_SHA=$(git rev-parse HEAD)
echo "📝 En son commit SHA: $NEW_SHA"

# 2. ACR Login Server
ACR_SERVER="veritasacr.azurecr.io"

# 3. Service listesi
SERVICES=("creator-service" "redirector-service" "analytics-service")

echo "🔄 Image tag'leri güncelleniyor..."

# 4. Her service için YAML dosyasındaki image tag'ini güncelle
for service in "${SERVICES[@]}"; do
    yaml_file="k8s/${service}.yaml"
    
    if [ -f "$yaml_file" ]; then
        # Eski image tag'ini yeni SHA ile değiştir
        sed -i "s|image: ${ACR_SERVER}/veritas/${service}:.*|image: ${ACR_SERVER}/veritas/${service}:${NEW_SHA}|g" "$yaml_file"
        echo "✅ ${service} YAML dosyası güncellendi"
    else
        echo "⚠️  ${yaml_file} dosyası bulunamadı, atlıyor..."
    fi
done

echo "🔄 Kubernetes'e deploy ediliyor..."

# 5. Kubernetes'e apply et
for service in "${SERVICES[@]}"; do
    yaml_file="k8s/${service}.yaml"
    
    if [ -f "$yaml_file" ]; then
        echo "📦 ${service} deploy ediliyor..."
        kubectl apply -f "$yaml_file"
    fi
done

echo "🎉 Deployment tamamlandı!"
echo "📊 Pod durumlarını kontrol et:"
echo "   kubectl get pods"
echo "   kubectl logs -l app=creator-service" 