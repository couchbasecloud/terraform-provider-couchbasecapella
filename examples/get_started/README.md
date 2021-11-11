# Example - Getting started with Couchbase Capella and Terraform

This example will cover setting up Terraform with Couchbase Capella. This will create the following resources in Couchbase Capella:

- Project
- Couchbase Cluster - m5.xlarge
- Bucket
- Database User
<!-- - IP Allowlist -->

<!-- ## Dependencies

- Terraform v0.13 or greater
- A Couchbase Capella account -->

## Usage

**1\. Setting up Authentication**

You will need to provide your credentials for authentication via the environment variables:

```bash
export CBC_ACCESS_KEY="xxxx"
export CBC_SECRET_KEY="xxxx"
```

**2\. Review the Terraform plan**

Execute the following command to review the resources that will be deployed.

```bash
$ terraform plan
```

This project currently creates the below deployments:

- Project
- Couchbase Cluster - m5.xlarge
- Bucket
- Database User
<!-- - IP Allowlist -->

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
