#! /usr/bin/env bash

kubectl apply -f crds/MarkhorSecret_crd.yaml
kubectl apply -f namespace.yaml
kubectl apply -f role.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rolebinding.yaml
kubectl apply -f secret.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
cat ./validation_webhook.yaml | sed "s/_CA_BUNDLE_HERE_/${CA_BUNDLE}/" | kubectl apply -f -

./gen_hook_tls_secret.sh
