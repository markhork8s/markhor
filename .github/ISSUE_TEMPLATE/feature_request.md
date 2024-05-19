---
name: Feature request 🐙
about: Help us improve
title: ""
labels: "enhancement"
assignees: "civts"
---

# Feature request 🐙

## How to use this template

- [ ] Complete all the sections of this template.
- [ ] Remove all the content inside the \<description\> tag. This is only to explain each section of the template.
- [ ] Make sure that no personal or sensitive information is present in this document before submitting it. If there is, replace it with asterisks "\*" or just replace it with generic values (e.g., replace your domain name with "example.com").
- [ ] **Delete this section** before submitting the issue.

## Pre-flight

<description style="display:none">
Check **all** the boxes before proceeding. To check a box, write it like this: [x]

</description>

- [ ] I checked that this is not a duplicate of existing issues/pull requests.

- [ ] I am quite sure this is a feature request for the Markhor project itself and not for any of the other technologies it interacts with (SOPS, Kubernetes, etc.).

- [ ] I tried to achieve the same result with alternative solutions and think that the best way to proceed is to add this functionality to Markhor. In my opinion, the effort to add, test and maintain this functionality is justified by the ease of use it brings.

- [ ] I have read and accepted the [Contributing Guidelines](https://github.com/markhork8s/markhor/blob/main/CONTRIBUTING.md) and the [Code of Conduct](https://github.com/markhork8s/markhor/blob/main/CODE_OF_CONDUCT.md).

## Explain your proposed feature/enhancement

### What is it?

[This is an example feature request, replace it with yours]  
Markhor should be able to send alerts with [gotify](https://gotify.net/) when a `MarkhorSecret` fails to decrypt.

### Why is it needed?

This is needed so that cluster administrators can know if their developers are having trouble adopting Markhor.

### How should it work?

We could use package [xyz](https://pkg.go.dev/google.golang.org/grpc/examples/helloworld/helloworld) for the implementation, and add a new top-level option called `notify` in the Markhor configuration file with x,y,z as sub.options

### Are there any existing workarounds or alternative solutions? Why are they insufficient?

It is possible to achieve the same by monitoring the logs of the program. [Here](example.com) is an example of how to do it. In this case, we wrote a program that watches the log file (specified with the `logging.additionalLogFiles` configuration option in Markhor) and sends an alert with gotify when there is an error.

I think that if this functionality was integrated directly in the main program, it would be easier for cluster administrators to set it up.

### How will it impact existing users?

By default, the option should be inactive, so there will be no difference for existing users.

### Would you be willing to contribute with the necessary code and tests?

Tell us if you'd be willing to contribute with code and tests for this feature, should the need arise.
