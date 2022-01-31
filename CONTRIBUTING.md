# Contributing

Couchbase welcomes and encourages community contributions to our Capella provider! Below are the following steps for making a contribution to this provider:

## Prerequisites

### Enable MFA

We require all contributors have [MFA enabled](https://help.github.com/en/github/authenticating-to-github/configuring-two-factor-authentication) on their account.

### Verify your commits

We require that contributors verify their commits. You’ll need to do this before GitHub will accept pushes for all branches.
Please follow the steps below:

1. [Generate a GPG Key](https://help.github.com/en/github/authenticating-to-github/generating-a-new-gpg-key)
2. [Add your public key to your Github Account](https://help.github.com/en/github/authenticating-to-github/adding-a-new-gpg-key-to-your-github-account)
3. [Tell your local git about your key](https://help.github.com/en/github/authenticating-to-github/telling-git-about-your-signing-key)

```
$ git config user.signingkey YOUR_KEY
```

4. [Configure git to sign your commits](https://help.github.com/en/github/authenticating-to-github/signing-commits) for the local repository

```
$ git config commit.gpgsign true
```

5. Make sure your [local git is configured to use the same email address](https://help.github.com/en/github/setting-up-and-managing-your-github-user-account/setting-your-commit-email-address) as your GPG key.

```
$ git config user.email "your.name@email.com"
```

## Before Contributing

<!-- Sign the contributing agreement -->

1. Please read through the [Terraform Contribution Guidelines](https://www.terraform.io/docs/extend/community/contributing.html) and the [README](https://github.com/couchbasecloud/terraform-provider-couchbasecapella/blob/main/README.md) in this repository.
2. Please [file an issue](https://github.com/couchbasecloud/terraform-provider-couchbasecapella/issues) for the contribution you'd like to make, where we can then discuss.

## After Issue Discussion

Once we've discussed your proposed contribution, and you're ready to implement your changes, please follow the following steps:

1. Get the latest changes from upstream `git checkout main` , `git pull`.
2. Create a branch with a name that describes the issue you filed.
3. Make sure your implementation adheres to [Terraform Best Practices](https://www.terraform.io/plugin/sdkv2/best-practices).
4. Please familiarise yourself with the [Terraform testing documentation](https://www.terraform.io/plugin/sdkv2/testing).
5. Commit and push your code to your branch. When writing your commit message keep the **first line of your commit message to 50 characters** or less, followed by a blank line, followed by an **explanation of the commit wrapped to 72 characters**.
6. Create a pull request.
7. We will then review your pull request, and may request additional changes or merge your pull request.
