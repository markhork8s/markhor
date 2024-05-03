<p align="center">
  <picture height="200px">
    <source media="(prefers-color-scheme: dark)" srcset=".github/images/logo/dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/images/logo/classic.svg">
    <img height="200px" alt="Markhor logo" src=".github/images/logo/classic.png">
  </picture>
</p>

# Markhor 🐐

-In its current state, this is a 100% experimental project. I do not know if this will ever work. DO NOT USE THIS CODE-

This is a program meant to run as a pod in your k8s cluster. Once installed -with the sops private key mounted in `/root/.config/sops/age/keys.txt`-, it watches for MarkhorSecret resources. When one is created, the pod will decrypt its content and create corresponding Kubernetes Secret. If you update/delete the MarkhorSecret, the generated secret will reflect the changes as well.

Check the Example usage section to see how to define the secrets.

This project is in no way affiliated or endorsed by SOPS nor Kubernetes.

# Table of Contents
# Introduction
# Setup
# Example usage
## Setup on the cluster
1. Install markhor in your cluster -along with the CRDs-
1. Create a serviceaccount, a role and a role-binding to give it permission to view the MarkhorSecret resources and manage the Secret ones.
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
- https://blog.pelo.tech/k-markhor-argocd-54becd3a1a34
- https://community.ops.io/jilgue/secrets-in-argocd-with-sops-pa6
- https://medium.com/@CoyoteLeo/security-upgrade-with-sops-5d4a1385c680
- https://www.redhat.com/en/blog/a-guide-to-gitops-and-secret-management-with-argocd-operator-and-sops

So I decided to adopt the same philosophy as reflector, which clones the secrets.

However, since the 'Secret' resource is defined by kubernetes, it is not possible to create a Secret that also has the `sops` property (TODO: add an explanation where you show what happens when you try to apply a secret encrypted with sops that fails because the base64 encoded data is incorrect and there is the additional property `sops`). This is the reason why I created a CRD for a MarkhorSecret. It may have been possible to extend the definition of a k8s secret, but this seems to me to have too much disruptive potential.

# Security concerns
## Confused deputy attack
Should not be possible since, thanks to the MAC (Message Authentication Code) that SOPS includes in its files, it is not possible to alter their content -even the parts which are unencrypted-.

# Limitations:
Due to how k8s marshals the applied configurations, comments will not persist in the final output.

Also, for the decryption to work, all the fields must be in the same order as the original file -otherwise, the MAC check fails-.

Also, no empty lines -except for the last one-

## Possible solutions:
- Optional field order in the Markhorsecret itself, since arrays are not reordered
  ```yaml
  dataOrder:
    - z/key.pem
    - certificate
  stringDataOrder:
    - session
  ```
- Another CRD: customorder with a spec field with an array of objects. The program, runtime, will load them

## TODOs:

- indentation does not count!
- comments do not count!
- it seems that all that is needed is for all the fields to be in the same order
- the order inside the 'sops' field does not matter
- it does not even need to be in the same format. You can encrypt as YAML, convert to JSON and decrypt! Provided the order rules are respected

# Config:

Markhor can get its configuration form one or more of these sources. Here they are listed in order of decreasing priority.

1. Environment variables
1. Config file
1. Runtime arguments

If a configuration value is specified in more than one source, a warning will be issued unless `logging.warn_on_value_override` is set to `false`.

## Values:

```yaml

kubeconfigPath: null #Path to the kubeconfig file to connect with the kubernetes cluster
kubernetesClusterTimeoutSeconds: 10
sopsKeysPath: ~/.config/sops/keys #Path to the keys that SOPS will use for decryption

healthcheck:
  port: 8080

logging:
  level: "info"

behavior:
  fieldmanager:  # The field manager for the k8s Secrets managed by Markhor. See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
    name: "github.com/civts/markhor"
    forceUpdates: false
  overrideExistingSecrets: false
  pruneDanglingSecrets: false #Deletes Secrets that have the managed-by Markhor annotation but no corresponding Markhor Secret

markorSecrets:
  hierarchySeparator:
    default: "/"
    allowOverride: false
    warnOnOverride: true
  managedAnnotation:
    default: "markhor.example.com/managed-by"
    allowOverride: false
    warnOnOverride: true
  namespaces: # Watch for Markhor secrets only in these namespaces
    - a
    - b

```

```
--config, -c   path to the config file (default /etc/markhor/config.yaml)
```

# Features:

1. Validation hook that checks if it is possible to decrypt a MarkhorSecret and if not tells you why -preventing its creation and updating-. ✅

1. Create a MS ✅
  1. If it can be decrypted, it does a kubectl apply, meaning it creates the Secret if it's missing, updates it if needed and leaves it unchanged if everything is the same
1. Modify a MS ✅
  Same as create since kubectl handles the apply
1. Delete a MS ✅
  The corresponding Secret is deleted

1. Healthcheck ✅

1. Logging with slog and different levels ✅

1. Honoring all the config values ❌

1. Docs ❌
  1. https://www.bestpractices.dev/en
  1. https://goreportcard.com/

1. Testing ❌
