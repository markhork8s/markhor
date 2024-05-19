# Markhor Security Policy 🏰

## Reporting Security Vulnerabilities

If you discover a security vulnerability in Markhor, please report it to us without delay. We take security vulnerabilities seriously and appreciate your help in responsibly disclosing any issues you may find.  
If you are unsure if what you found is a vulnerability, do treat it as such.

To report a security vulnerability, please visit: [https://github.com/markhork8s/markhor/security/advisories](https://github.com/markhork8s/markhor/security/advisories). If for any reason that is inaccessible, use the contacts you find [here](https://zanol.li/#contacts), but only the ones that provide an encrypted channel -i.e., matrix-. Do not use any public or unencrypted channel to communicate with us about security vulnerabilities.

In your report, please include the following:

1. A description of the vulnerability
1. Steps to reproduce it
1. Potential impacts and risks
1. Which version of Markhor, K8s, SOPS, etc. you used
1. Any proposals you may have on how to fix it

## Disclosure Policy

If you are eager to share with others your discovery, we understand the enthusiasm and will be there to celebrate with you once it is patched.  
At the same time, we need time to properly assess and address the issue. Therefore, we request that you _do not_ share the vulnerability with anyone other than the project's maintainers until we have publicly released a fix OR 90 days have passed since when you reported it to us.

We will recognize in the release notes that you found that vulnerability (and if you contributed to the fix).

## Maintainer Commitment

As the current sole maintainer of Markhor, I am committed to addressing any security issues that may arise to the best of my ability. Please understand that as a solo maintainer, I may not be able to provide immediate fixes, but I will work diligently to resolve any security concerns in a timely manner.

That said, any help is appreciated since I'd like this project to grow and be community-driven.

## Threat Model

The security of this program relies on the following elements:

- the security of [SOPS](https://github.com/getsops/sops)
- the security of the encryption method chosen (age, gpg, etc.)
- the correct management of the private key(s) used for encryption and decryption (it should be available only to the Markhor pod)
- the correct management of the Kubernetes cluster. Users and other pods must not be able to:
  - read the Secret containing the decryption key
  - create, read, modify or remove `Secret`s that are managed by Markhor
- having an honest cluster administrator (if they are the ones to go rogue, game over for us)
- the security of the golang dependencies (we have added SAST checks in the CI, and dependabot)

During the development of this project, we are considering the following threats:

1. **Unauthorized CRUD operations on resources**

Markhor allows only to act on `Secret`s. Markhor maintains a one-to-one correspondence between the `MarkhorSecret`s and the relative `Secret`s. Each managed `Secret` is created with the same name and namespace as the `MarkhorSecret`. Therefore, as long as you apply the same policies that you use to regulate who can create `Secret`s where to `MarkhorSecret`s, malicious users will not be able to tamper with `Secret`s they should not access. If some of your users only use Markhor to manage the `Secret`s, they do not need to be able to interact with `Secret`s directly anymore and that authorization should be revoked.

1. **Tampering of the Secret data**

Markhor protects the integrity and the confidentiality of the `MarkhorSecret`s using SOPS. Once a file has been cyphered with SOPS -assuming there are no vulnerabilities in SOPS-, it can only be decrypted with the original key and the decryption will fail if the data has been altered -we use the decryption function provided by SOPS, which checks the MAC-.

1. **Denial of Service attack**

It should not be possible to crash the Markhor application by submitting invalid data.
It also should not be possible to alter the content of an existing `Secret` with invalid data -thanks to the MAC-.

The adversary we consider is a rogue developer. They have the following capabilities:

1. Can read any `MarkhorSecret`, but only encrypted (otherwise the problem is with the SOPS workflow the other devs use)
1. Can tamper any encrypted `MarkhorSecret` (but the corresponding `Secret` will not be altered since the MAC breaks)
1. Can delete a `MarkhorSecret` (in this case, the corresponding `Secret` gets deleted)
