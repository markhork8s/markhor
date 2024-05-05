#! /usr/bin/env bash

set -e

kubectl apply -f crds/MarkhorSecret_crd.yaml
kubectl apply -f namespace.yaml
kubectl apply -f role.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rolebinding.yaml
kubectl apply -f private_key_secret.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml

# This section is necessary only if you use the admission controller
kubectl apply -f service.yaml
CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
cat ./validation_webhook.yaml | sed "s/_CA_BUNDLE_HERE_/${CA_BUNDLE}/" | kubectl apply -f -
echo "Remember to run ./gen_hook_tls_secret.sh if you need them for the admission controller" 
