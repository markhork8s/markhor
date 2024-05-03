#!/usr/bin/env bash

# Get program version
V=$(cat pkg/costants.go | grep " VERSION " | sed 's/.*= "//' | sed 's/"//')
IMAGE=markhor/markhor:$V

# Build image
docker build -t $IMAGE .

# Prepare for import in k8s
sed -i "s/markhor\/markhor:.*/markhor\/markhor:$V/" helm/deployment.yaml
minikube image load $IMAGE
# docker save $IMAGE | sudo k3s ctr images import -

