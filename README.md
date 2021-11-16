# Couchbase Capella Provider

This is the repository for the Terraform Couchbase Capella Provider which allows you to use Terraform with Couchbase Capella.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
- [Go](https://golang.org/doc/install) >= 1.14

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `build` command:

```sh
$ go build
```

After the provider has been built locally it must be placed in the user plugins directory so it can be discovered by the
Terraform CLI. Please execute the following command to move the provider binary to this directory:

```sh
$ mv terraform-provider-couchbasecloud ~/.terraform.d/plugins/terraform.couchbase.com/local/couchbasecloud/1.0.0/<OS_ARCH>
```

The terraform provider is installed and can now be discovered by Terraform through the following HCL block.

```hcl
terraform {
  required_providers {
    couchbasecapella = {
      source  = "terraform.couchbase.com/local/couchbasecloud"
      version = "1.0.0"
    }
  }
}
```

## Using the provider

### Configuring Programmatic Access

In order to set up authentication with the Couchbase Capella provider a programmatic API key must be generated. Instructions to generate your API key can be found in the [Couchbase Capella Public API documentation](https://docs.couchbase.com/cloud/public-api-guide/using-cloud-public-api.html).

### Authenticating the Provider

You will need to provide your credentials for authentication via the environment variables,
`CBC_ACCESS_KEY` and `CBC_SECRET_KEY`,
for your access and secret API Key Pair respectively.

Usage (prefix the export commands with a space to avoid the keys being recorded in OS history):

```shell
$  export CBC_ACCESS_KEY="xxxx"
$  export CBC_SECRET_KEY="xxxx"
$ terraform plan
```

### Example Usage

```terraform
#Configure the Couchbase Capella Provider
provider "couchbasecapella" {}

# Create example project resource
resource "couchbasecapella_project" "project" {
  name = "project1"
}
```

**1\. Initialise the Terraform provider**

Execute the following command to initialise the terraform provider.

```bash
$ terraform init
```

**2\. Review the Terraform plan**

Execute the following command to review the resources that will be deployed.

```bash
$ terraform plan
```

**3\. Execute the Terraform apply**

Execute the plan to deploy the Couchbase Capella resources.

```bash
$ terraform apply
```

**4\. Destroy the resources**

Execute the following command to destroy the resources so you avoid unnecessary charges.

```bash
$ terraform destroy
```

### Getting Started

Please also visit the `get_started` directory for an example configuration for provisioning a project, cluster, bucket and database user.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above). You will also need to have access to a [Couchbase Capella](https://www.couchbase.com/products/capella) account.

## Testing the Provider

### Configuring the environment variables

You must also configure the following environment variables before running the tests:

```sh
export CBC_CLOUD_ID=<YOUR_CLOUD_ID>
export CBC_PROJECT_ID=<YOUR_PROJECT_ID>
export CBC_CLUSTER_ID=<YOUR_CLUSTER_ID>
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
