# Couchbase Capella Provider

You can use the Couchbase Capella provider to interact with Projects, Clusters, Buckets and Database Users within your Couchbase Capella tenant.

The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available provider resources and data sources.

## Example Usage

```terraform
#Configure the Couchbase Capella Provider
provider "couchbasecapella" {}
# Create the resources
```

## Configuring Programmatic Access

In order to set up authentication with the Couchbase Capella provider a programmatic API key must be generated. Instructions to generate your API key can be found in the [Couchbase Capella Public API documentation](https://docs.couchbase.com/cloud/public-api-guide/using-cloud-public-api.html).

## Authenticating the Provider

You will need to provide your credentials for authentication via the environment variables,
`CBC_ACCESS_KEY` and `CBC_SECRET_KEY`,
for your access and secret API Key Pair respectively.

Usage (prefix the export commands with a space to avoid the keys being recorded in OS history):

```shell
$  export CBC_ACCESS_KEY="xxxx"
$  export CBC_SECRET_KEY="xxxx"
$ terraform plan
```
