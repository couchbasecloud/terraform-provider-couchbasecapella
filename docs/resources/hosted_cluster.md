---
page_title: "Couchbase Capella: Hosted Cluster"
subcategory: ""
description: |-
Create, edit and delete Hosted Clusters in Couchbase Capella.
---

# Resource couchbasecapella_hosted_cluster

`couchbasecapella_hosted_cluster` allows you to create, edit and delete hosted clusters in Couchbase Capella. The resource requires your Project ID.

~> **WARNING:** This current release of Terraform Couchbase Capella provider doesn't support creating bucket or database user resources for hosted clusters. Please log in to the Couchbase Capella UI where you'll be able to manage buckets and database users once your hosted cluster has been deployed.

~> **WARNING:** **UPDATING CLUSTER SERVERS WILL CAUSE BUCKETS TO BE DELETED**. Updating cluster servers will cause your cluster to redeploy. You won't be able to access the cluster in the Couchbase Capella UI until it has been deployed.

## Example Usage

### Example Hosted Cluster

```hcl
resource "couchbasecapella_hosted_cluster" "test" {
  name        = "cluster_name"
  project_id  = "your_project_id"
  place {
    single_az = true
    hosted {
      provider = "aws"
      region   = "us-west-2"
      cidr     = "cidr_block"
    }
  }
  support_package {
    timezone = "GMT"
    support_package_type     = "Basic"
  }
  servers {
    size     = 3
    compute  = "m5.xlarge"
    services = ["data"]
    storage {
      storage_type = "GP3"
      iops = "3000"
      storage_size = "50"
    }
  }
}
```

## Argument Reference

- `name` - (Required) The name of the cluster you want to create. The cluster name can include letters, numbers, spaces, periods (.), dashes (-), and underscores (\_). Cluster name should be between 2 and 128 characters and must begin with a letter or a number.
- `project_id` - (Required) The id of the project where your cluster will be created. This must be a valid UUID and an existing project ID. (Cannot be changed via this Provider after creation.)
- `description` - (Optional) A description for the cluster.

### Place

- `single_az` - (Required) A boolean to describe if there is only a single availability zone. (Cannot be changed via this Provider after creation.)

~> **WARNING:** `single_az` must be set to true if the you select the "Basic" support package.

#### Hosted

- `provider` - (Required) The name of the cloud provider you want your cluster to be hosted in. `aws`, `azure`, or `gcp` are the available providers that you can specify. (Cannot be changed via this Provider after creation.)
- `region` - (Required) A valid region for the cloud provider that you want you cluster to be hosted in. This must be a valid region for the cloud provider you have specified. (Cannot be changed via this Provider after creation.)
- `cidr` - (Required) The CIDR address. This must be a valid CIDR address. (Cannot be changed via this Provider after creation.)

### Support Package

- `timezone` - (Required) The time zone that you would like to receive support from. `ET`, Eastern Time, `GMT`, Greenwich Mean Time , `IST`, India Standard Time, `PT`, Pacific Time, are the available time zones that you can specify.
- `support_package_type` - (Required) The support plan that you would like for your Capella cluster. `Basic`, `DeveloperPro`, `Enterprise` are the available support plans that you can specify.

For more detailed information on support packages, you can view this [detailed plan comparison](https://www.couchbase.com/support-policy/cloud).

### Servers

- `size` - (Required) The number of nodes in your cluster. This must be a value between 3 and 27.
- `compute` - (Required) The name of the compute instance type. This must be a valid compute instance type for the provider that you have specified.
- `services` - (Required) A list of Couchbase services that you want in your cluster. `Data`, `Query`, `Index`, `Search`,`eventing`, `analytics` are the available services that you can specify.

#### Storage

- `storage_type` - (Required) The name of the storage type. `GP3`, `IO2` are the available storage types that you can specify.
- `iops` - (Required) The number of the IOPS.

~> **IMPORTANT:** The minimum value is 3000 for GP3, and 1000 if IO2. The maximum value is 16000 for GP3, and 64000 if IO2.

- `storage_size` - (Required) The storage per node in gigabytes.

~> **IMPORTANT:** The minimum storage per node is 50Gb. The maximum storage per node is 16Tb.

## Attribute Reference

- `id` - The cluster id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clustersv3).
