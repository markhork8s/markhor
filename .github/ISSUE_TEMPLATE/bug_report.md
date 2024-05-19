---
name: Bug report 🐞
about: Help us improve
title: ""
labels: "bug"
assignees: "civts"
---

# Bug report 🐞

## How to use this template

- [ ] Complete all the sections of this template.
  - When you see "<output of `some_command`>", run in the terminal the command `some_command` and copy its output. Quote the resulting text with the \` symbol to denote it as code.
- [ ] Remove all the content inside the \<description\> tag. This is only to explain each section of the template.
- [ ] Make sure that no personal or sensitive information is present in this document before submitting it. If there is, replace it with asterisks "\*" or just replace it with generic values (e.g., replace your domain name with "example.com").
- [ ] **Delete this section** before submitting the issue.

## Pre-flight

<description style="display:none">
Check **all** the boxes before proceeding. To check a box, write it like this: [x]

</description>

- [ ] I checked that this is not a duplicate of existing issues/pull requests.

- [ ] I am quite sure this is an issue with the Markhor project itself and not with any of the other technologies it interacts with (SOPS, Kubernetes, etc.).

- [ ] I spent time to ensure others can reproduce my issue as well.

- [ ] I have read and accepted the [Contributing Guidelines](https://github.com/markhork8s/markhor/blob/main/CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md).

## Expected behavior

<description style="display:none">
In this section, write about: 
- What you were trying to do
- How you tryied to do it
- What you expected to see as a result
</description>

## Actual behavior

<description style="display:none">
This section should answer the following question:
What did you see instead of the expected result?
</description>

## Detailed steps to reproduce

<description style="display:none">
Include here relevant Kubernetes manifests, command outputs, error messages or logs. If the previous sections already contain everything, you can delete this one.

Example:

1. Clone the repo
1. Run this command `echo "abcd"`
1. Apply the manifest for service foo:
   ```yaml
   apiVersion: foo
   kind: bar
   ```
1. You should get the error in the logs: `could not do X because Y failed at Z`

</description>

## Your environment setup

1. Golang version: <output of `go version`>
1. Linux version: <output of `uname -s -r -v -p -i -o`>
1. Kubernetes version: <output of `kubectl version --short`>
1. SOPS version: <output of `sops --version`>

## If the bug is confirmed, would you be willing to submit a pull request for this?

- [ ] Yes
- [ ] No
