# Markhor Contributing Guidelines

Looking for how you can lend a hand? Glad to have you here! 🔥🎆

Wether you want to write 0, 20 or 2000 lines of code to help this project grow, this is the right place to start!

Thank you for taking time to consider contributing to Markhor ❤️

# Introduction

### Why you should read this document

> Following these guidelines shows that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

### How you can contribute

Here are some things you can do to help the Markhor project:

- **Use it** in your activities and report any issues you find so we can fix them for everyone 🧪
- Improve the clarity, fluency or completeness of the **documentation** 📝
- **Share the project** with others who may be interested. Do it in person, with posts or videos on social medias, with blog posts. Wherever you want! Please, make sure to include a link back ot this repo whenever possible 📯
- **Propose** new features in the issues or **vote** for existing ones 🙋
- **Validate** existing issues or bugs by reproducing them and/or **test** changes made in pull requests 🦬
- Contribute **code** for either implementation of new features, maintenance or expanding the test suite 👽

These are only suggestions. Chances are you will find ways to contribute we didn't even expect. Keep an open mind!

# Ground Rules

### Communication Style

This project follows the Contributor Covenant code of conduct. We strive to honor it in every interaction and ask you to do the same.

For all our communications, we use English as the only language.

### Committment

The author, the mainteners and the contributors of this project devote their time to it mostly for passion.

We welcome everyone to contribute, but ask them to ensure that their work on this project does not become too great of a source of stress (we suggest contributors to read this [article](https://opensource.guide/maintaining-balance-for-open-source-maintainers/) on how to manage it).

None of the contributors is under any obbligation. This includes:

- reviewing comments/issues/pull requests 🗃
- replying within specific time-frames 🐇
- helping you debug issues with the software 🔍

This is not to say that we will not do these things. Just, do not take these for granted and show understanding and gratitude when appropriate.

If you find that the turnaround time for your requests is getting too long, consider supporting the project more actively 😜

# Getting started

## Dev environment setup

The suggested development environment for this project is by using the included nix flake. This way we can ensure everyone has the exact same setup. If you have no idea what a nix flake is, no worries. Just ignore this paragraph and use whatever development setup you are most comfortable with.

## Your First Contribution

If it is one of your first times contributing to open source projects, I strongly sugget you to take 30 minutes to read https://opensource.guide/how-to-contribute/. It is a curated series of articles that makes you aware of many little details that make great differences for both you and the maintainers of this project.

If you want to contribute by solving an existing issue, look for the ones with the tag `good first issue`. These are beginner issues and should take less than 100 lines of code to complete, plus the relative tests.
If you are more experienced, `help wanted` issues may be the ones for you.

## How to submit a contribution

Before starting to code, please review existing issues and pull requests to see if anyome has been working on something similar.
Also, note that we aim to have a comprehensive test coverage of the project, so your contribution should include tests -see the dedicated [section](#tests)-.

The recommended workflow is the following:

1. Create your own fork of the code
2. Do the changes in your fork
3. If you like the changes and think the project could use them:
   - Be sure you have added the tests too, that they pass and that any previous test passes too
   - Create the pull request

## Tests

Tests help ensuring the code is correct and preventing regressions over time.
The tests should not modify the state of the machine they run on and they should work offline. If a dependency is needed -like a kubernetes cluster or a REST server-, it has to be simulated (aka, mocked). It is acceptable to create temporary files on the machine running the tests using functions like mktemp.

You should always test the components you modify in isolation (unit tests).
In unit test, you will generally want to test that:

- the code behaves as expected when the input is correct
- the code behaves as expected (throwing an error) when the input is incorrect. This second case may include inputs in the wrong format, with missing data, with data out of range (for example, a negative UNIX timestamp)

You should also test how your component interacts with others using integration tests. In these tests, you will want to test both using correct and incorrect inputs.

# Issues and feature requests

Apart from security issues, issues, bugs, feature requests and enhancement proposals are all tracked on the repository's [issue page](https://github.com/markhork8s/markhor/issues/).

## Security-related issues

> If you found something that you think can compromise the **security** of the application or its users, go to the document [SECURITY.md](https://github.com/markhork8s/markhor/blob/main/SECURITY.md) instead.
>
> In order to determine whether you are dealing with a security issue, ask yourself these two questions:
>
> - Can I access something that's not mine, or something I shouldn't have access to?
> - Can I disable something for other people?
>
> If the answer to any of those two questions are "yes", then you're probably dealing with a security issue. Note that even if you answer "no" to both questions, you may still be dealing with a security issue, so if you're unsure, just treat your issue as such.  
> (this section was adapted from: [Travis CI](https://github.com/travis-ci/travis-ci/blob/master/CONTRIBUTING.md))

## Guidelines for discussing issues

When engaging with existing issues in the project, please follow these guidelines to ensure effective communication and collaboration:

1. **Upvoting Issues**:

   If you find an existing issue that resonates with you, please add a reaction to the opening comment to upvote the issue. We encourage using the 👍 emoji as a way to indicate your interest and support. This action is equivalent to a "+1".

2. **Example of Effective Commenting**:

   Comments that contribute valuable information about reproducing the bug or identifying affected versions are highly appreciated.

   For instance, if the issue is related to a problem with building the project's container image using Docker version `24.0.5`, a good comment could provide additional context, such as other affected versions:

   - "Hello everyone, I encountered a similar issue using Docker version `22.0.3`, but I saw that with `21.3.2` the image builds as expected."

   Or give additional info on the probable cause of the issue:

   - "This should be related to the change xyz that Docker did here (link to the PR/changelog in Docker when the change was introduced) where they changed how foo parses bar"

3. **Comments to Avoid**:

   Please refrain from posting comments that simply state "+1" or "same here." Instead, express your support by using the reaction feature.

   Avoid excessive pinging of maintainers. If you notice that a maintainer has already been notified, kindly wait for a reasonable period, such as a couple of weeks, before reaching out again.

## How to open an issue / feature request

1. Search the [issues on GitHub](https://github.com/markhork8s/markhor/issues/) to understand if anyone else had the same issue. Keep in mind that they may have described the problem using different terminology than you.  
   If you see that the issue has been inactive for more than a couple of months (i.e., it is stale), you may add a comment to revive it.
1. If nobody had filed a similar issue, proceed in opening a new issue [here](https://github.com/markhork8s/markhor/issues/new) (be sure to follow the templates for new issues).
1. Be prepared to engage in discussions with maintainers and other contributors. Try to respond timely to any requests for additional information or clarifications to expedite the resolution process.

# Community

As of now there are no community channels. No discord, no matrix, nada.

# Coding style

We adopted `gofmt` to format the code. It ensures a consistent coding style.

Try to write clean code, with descriptive yet concise variable names and functions that are not too convoluted. The code should be readable and maintainable.

# Git guidelines

## GPG signature

Do sign your commits with GPG to ensure they really come from you. See [this article](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) for the details on how this works on GitHub. We also suggest you test this functionality in a private dummy repository before using it here.

## Commit messages

Try to follow "The seven rules of a great Git commit message" ([link](https://cbea.ms/git-commit/)):

1. Separate subject from body with a blank line
1. Limit the subject line to 50 characters
1. Capitalize the subject line
1. Do not end the subject line with a period
1. Use the imperative mood in the subject line
1. Wrap the body at 72 characters
1. Use the body to explain _what_ and _why_ vs. how

---

This document has been adapted from [nayafia's contributing template](https://github.com/nayafia/contributing-template)
