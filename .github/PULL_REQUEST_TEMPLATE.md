## How to use this template
- [ ] Complete all the sections of this template.
- [ ] Remove all the content inside the \<description\> tag. This is only to explain each section of the template.
- [ ] Make sure that no personal or sensitive information is present in this document before submitting it. If there is, replace it with asterisks "*" or just replace it with generic values (e.g., replace your domain name with "example.com").
- [ ] Remove this section before submitting the PR.

## Pre-flight
<description style="display:none">
Check **all** the boxes before proceeding. To check a box, write it like this: [x]

</description>

- [ ] I checked that this is not a duplicate of existing issues/pull requests.

- [ ] This PR solves an existing issue (please link it) or is a trivial fix.

- [ ] I have read and accepted the [Contributing Guidelines](https://github.com/markhork8s/markhor/blob/main/CONTRIBUTING.md) and the [Code of Conduct](https://github.com/markhork8s/markhor/blob/main/CODE_OF_CONDUCT.md).

## Details

### What does this PR do?

[This is an example pull request, replace it with yours]  
Enables Markhor to send alerts with [gotify](https://gotify.net/) when a `MarkhorSecret` fails to decrypt.
It does so by using package xyz.
It also adds the Markhor configuration option `notify` and its relative entry in the JSON schema. 

### Does this PR fix any open issues?
This PR fixes issue #1234 (change the number of the issue appropriately)

### How will it impact existing users?

By default, the option will be inactive, so there will be no difference for existing users.

### Code quality

Did you
- [ ] Format the code with `gofmt`
- [ ] Write the tests for the functionality you added
- [ ] Ensure you did not break any of the pre-existing tests
