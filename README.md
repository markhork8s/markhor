# SOPS_K8S

-In its current state, this is a 100% experimental project. I do not know if this will ever work. DO NOT USE THIS CODE-

This is a program meant to run as a pod in your k8s cluster. Once installed -with the sops private key mounted in `/root/.config/sops/age/keys.txt`-, it watches for SopsSecret resources. When one is created, the pod will decrypt its content and create corresponding secret. If you update/delete the SopsSecret, the generated secret will reflect the changes as well.

Check the Example usage section to see how to define the secrets.

This project is in no way affiliated or endorsed by SOPS nor Kubernetes.

# Example usage
## Setup on the cluster
1. Install sops_k8s in your cluster -along with the CRDs-
1. Create a serviceaccount, a role and a role-binding to give it permission to view the SopsSecret resources and manage the Secret ones.
## Setup on the app
1. Create a `sops.yaml` file
1. Create a sops-encrypted secret. Ensure that no comments or empty lines are present in the file (see the limitations section).
1. If the ordering of the fields is not alphabetical, define a CustomOrdering in the same namespace.
1. kubectl apply the file(s) and check if the corresponding secrets were created.

# Why does this project exist?
I have been using SOPS for some time and so far it has worked great.
I wanted to use it to encrypt my k8s secrets so that I could keep them in my repo and deploy them with my CD.

While SOPS is supported for Flux, the support for ArgoCD is only available through plugins.
Now, looking at these resources, the integrations seem a bit too hacky for me
- https://blog.pelo.tech/k-sops-k8s-argocd-54becd3a1a34
- https://community.ops.io/jilgue/secrets-in-argocd-with-sops-pa6
- https://medium.com/@CoyoteLeo/security-upgrade-with-sops-5d4a1385c680
- https://www.redhat.com/en/blog/a-guide-to-gitops-and-secret-management-with-argocd-operator-and-sops

So I decided to adopt the same philosophy as reflector, which clones the secrets.

However, since the 'Secret' resource is defined by kubernetes, it is not possible to create a Secret that also has the `sops` property (TODO: add an explanation where you show what happens when you try to apply a secret encrypted with sops that fails because the base64 encoded data is incorrect and there is the additional property `sops`). This is the reason why I created a CRD for a SopsSecret. It may have been possible to extend the definition of a k8s secret, but this seems to me to have too much disruptive potential.

# Security concerns
## Confused deputy attack
Should not be possible since, thanks to the MAC (Message Authentication Code) that SOPS includes in its files, it is not possible to alter their content -even the parts which are unencrypted-. Also, the program creates the Seccret only in the same namespace where the SopsSecret was created.

# Limitations:
K8s sorts the fields alphabetically recursively and removes the comments when doing a kubectl apply.

For the decryption to work, all the fields and comments must be in the same order as the original file -otherwise, the MAC check fails-.

This means that the developers must write the yaml with the keys -and subkeys- sorted and without comments.

Also, no empty lines -except for the last one-

## Possible solutions:
- Optional field order in the sopssecret itself, since arrays are not reordered
  ```yaml
  dataOrder:
    - z/key.pem
    - certificate
  stringDataOrder:
    - session
  ```
- Another CRD: customorder with a spec field with an array of objects. The program, runtime, will load them

#### Custom client
https://pkg.go.dev/k8s.io/client-go/examples#section-readme
https://aimuke.github.io/k8s/2021/01/28/k8s-access-crds-from-client-go/
https://medium.com/cloud-native-daily/kubernetes-crd-handling-in-go-d426e9c3c1ab
https://medium.com/@disha.20.10/building-and-extending-kubernetes-a-writing-first-custom-controller-with-go-bc57a50d61f7
https://github.com/dishavirk/first-custom-k8s-controller
https://github.com/kubernetes/apiextensions-apiserver/tree/master/examples/client-go
https://github.com/kubernetes/client-go/tree/v0.29.3/examples/fake-client
https://github.com/dishavirk/canary-k8s-operator/tree/master

## TODOs:
- customordering CRD cant be recursive, it seems