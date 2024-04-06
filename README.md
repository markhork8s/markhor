# SOPS_K8S

-In its current state, this is a 100% experimental project. I do not know if this will ever work so do not use this code-

This is a simple program meant to run as a pod in your k8s cluster. Once installed -with the sops private key mounted in `/root/.age/keys.txt`-, create a secret encrypted with SOPS and add the annotation `sops_k8s/decryption-enabled: "true"`. The pod will decrypt the content and create a new secret in the same namespace. If you update/delete the original secret, the decrypted one will reflect the cahges as well.

This project is conceptually similar to refector. Only, it decrypts sops-encrypted values inside of secrets instead of replicating them.

# Why does this exist?
I have been using SOPS for some time and so far it has worked great.
I wanted to use it to secure my secrets so that I could keep them in my repo and deploy them with my CD.

While SOPS is supported for Flux, the support for ArgoCD is only available through plugins.
Now, looking at these resources the integrations seem a bit too hacky for me
- https://blog.pelo.tech/k-sops-k8s-argocd-54becd3a1a34
- https://community.ops.io/jilgue/secrets-in-argocd-with-sops-pa6
- https://medium.com/@CoyoteLeo/security-upgrade-with-sops-5d4a1385c680
- https://www.redhat.com/en/blog/a-guide-to-gitops-and-secret-management-with-argocd-operator-and-sops

So I decided to adopt the same philosophy as reflector.

# Security concerns
## Confused deputy attack
Should not be possible since, thanks to the MAC, I should not be able to take one value from a yaml and include it in another one,but I need to check this.
