To configure markhor to use a method xyz supported by SOPS, we use either env variables or files that SOPS will look for.
Therefore, mainly refer to the SOPS documentation. This document is to be considered a cheatsheet/quickstart.

# Age

1. Generate the private and public key pair: `age-keygen`
1. Create on the cluster a Secret with the private key
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: markhor-age-secret
     namespace: markhor
     labels:
       app.kubernetes.io/name: markhor
       app.kubernetes.io/instance: markhor-default
       app.kubernetes.io/version: 1.0.0
       app.kubernetes.io/component: operator
       app.kubernetes.io/part-of: markhor
   stringData:
     age_keys.txt: "AGE-SECRET-KEY-1LYQ3PW...."
   ```
1. Mount the secret in the markhor container by editing the deployment:
   ```yaml
   # This snippet contains only the relevant parts, it is not a complete deployment
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: markhor
   spec:
     spec:
       containers:
         - name: markhor
           env:
             - name: SOPS_AGE_KEY_FILE
               value: /age-secrets/age_keys.txt
           volumeMounts:
             - name: markhor-age
               mountPath: /age-secrets
               readOnly: true
       volumes:
         - name: markhor-age
           secret:
             secretName: markhor-age-secret
   ```
1. On the machine where you need to encrypt the MarkhorSecrets, create the `.sops.yaml` file
   ```yaml
   keys:
     - &mykey age1apq7ck5adq6dkd0c242phl42fsurvpxvt9pwk0qg7ahdex7fqppqj8pe8y
   creation_rules:
     - path_regex: ".*_secret.ya?ml"
       key_groups:
         - age:
             - *mykey
       encrypted_regex: ^(data|stringData)$
   ```

# AWS KMS

This is a basic setup, modify yours as needed

1. Go to your aws iam console, https://us-east-1.console.aws.amazon.com/iam/home
1. Click on "users"
1. Create user
1. Click on the user name of the newly creted user
1. Go to the tab "security credentials"
1. Cerate a file in `/.aws/credentials` on the machine where you need to encrypt the MarkhorSecrets
   ```
   [default]
   aws_access_key_id = AKI...
   aws_secret_access_key = ...
   ```
1. Cerate a Secret in the cluster where you need to decrypt the MarkhorSecrets
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: aws-kms-secret
     namespace: markhor
   stringData:
     access_key_id: AKI...
     secret_access_key: yoursecretaccesskey
   ```
1. Mount the secret in the markhor deployment populating the env of the container:
   ```yaml
   # This snippet contains only the relevant parts, it is not a complete deployment
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: markhor
   spec:
     spec:
       containers:
         - name: markhor
           env:
             - name: AWS_ACCESS_KEY_ID
               valueFrom:
                 secretKeyRef:
                   name: aws-kms-secret
                   key: access_key_id
             - name: AWS_SECRET_ACCESS_KEY
               valueFrom:
                 secretKeyRef:
                   name: aws-kms-secret
                   key: secret_access_key
       volumes:
         - name: aws-kms-secret
           secret:
             secretName: aws-kms-secret
   ```
1. Go to your aws kms console, https://us-east-1.console.aws.amazon.com/kms/home
1. Create a symmetric key
1. Click on its id in the table, copy the ARN (will look like this `arn:aws:kms:us-east-1:005113765608:key/727c4c86-6b04-4143-93a9-d4b269cfc7a5`)
1. Create the `.sops.yaml` file on the machine where you need to encrypt the MarkhorSecrets:
   ```yaml
   creation_rules:
     - kms: arn:aws:kms:us-east-1:005113765608:key/727c4c86-6b04-4143-93a9-d4b269cfc7a5
       encrypted_regex: ^(data|stringData)$
       ...
   ```

# GCP KMS

This example uses service accounts, which [may not be the most secure solution](https://cloud.google.com/docs/authentication/application-default-credentials). If you are more familiar with GCP than me, better methods are welcome.

1. Go to the [GCP KMS console](https://console.cloud.google.com/security/kms)
1. Create the keyring and the key using the cloud shell (I followed this [tutorial](https://codelabs.developers.google.com/codelabs/encrypt-and-decrypt-data-with-cloud-kms)),
1. Configure sops to use that key:
   ```yaml
   creation_rules:
     - gcp_kms: projects/MYPROJECT/locations/global/keyRings/MYKEYRING/cryptoKeys/MYKEYNAME
       encrypted_regex: ^(data|stringData)$
   ```
1. Create the service account credentials ([docs](https://cloud.google.com/iam/docs/keys-create-delete#creating))
   1. Go to https://console.cloud.google.com/apis/api/cloudkms.googleapis.com/credentials > credentials and create a new service account
   1. Assign the roles: `Cloud KMS CryptoKey Encrypter/Decrypter` and `Cloud KMS CryptoKey Public Key Viewer`
   1. Click the email address of the service account that you want to create a service account key for > Keys > add key > json
1. Set the env `GOOGLE_APPLICATION_CREDENTIALS` to the path of the JSON file on the machine you want to perform the encryption from.
1. In the cluster, create a Secret with the service account key, like this:
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: gcp-kms-secret
     namespace: markhor
   stringData:
     account_credentials.json: |
       {
         "type": "service_account",
         "project_id": "PROJECT_ID",
         "private_key_id": "KEY_ID",
         "private_key": "-----BEGIN PRIVATE KEY-----\nPRIVATE_KEY\n-----END PRIVATE KEY-----\n",
         "client_email": "SERVICE_ACCOUNT_EMAIL",
         "client_id": "CLIENT_ID",
         "auth_uri": "https://accounts.google.com/o/oauth2/auth",
         "token_uri": "https://accounts.google.com/o/oauth2/token",
         "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
         "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/SERVICE_ACCOUNT_EMAIL"
       }
   ```
1. Edit the deployment to use the Secret with teh GCP credentials:
   ```yaml
   # This snippet contains only the relevant parts, it is not a complete deployment
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: markhor
   spec:
     spec:
       containers:
         - name: markhor
           env:
             - name: GOOGLE_APPLICATION_CREDENTIALS
               value: /gcp/account_credentials.json
           volumeMounts:
             - name: gcp-kms
               mountPath: /gcp
               readOnly: true
       volumes:
         - name: gcp-kms
           secret:
             secretName: gcp-kms-secret
   ```
