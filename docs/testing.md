# Minikube

1. `nix-shell -p minikube kubectl`
   To get minikube and kubectl
1. `minikube start`
   To start the cluster
1. `minikube ssh -- "sudo cat /etc/kubernetes/admin.conf" | sed "s/control-plane.minikube.internal/$(minikube ip)/" > cluster.conf`
   1. `minikube ssh -- "sudo cat /etc/kubernetes/admin.conf" > cluster.conf`
      To create `cluster.conf` in the current directory with the information necessary to connect with the cluster.
   1. In `cluster.conf`, replace 'control-plane.minikube.internal' with the output of `minikube ip`
1. Open VSCode
1. Ensure the path to `cluster.conf` is correct in `.vscode/launch.json`
1. Press F5

# K3S

`export KUBECONFIG=/etc/rancher/k3s/k3s.yaml`
