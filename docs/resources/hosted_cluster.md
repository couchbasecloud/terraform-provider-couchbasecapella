# couchbasecapella_hosted_cluster Resource

`couchbasecapella_hosted_cluster` provides a Hosted Cluster resource. The resource allows you to create, edit and delete hosted clusters. The resource requires your Project ID.

~> **WARNING:** This current release of Terraform Couchbase Capella provider doesn't support creating bucket or database user resources for hosted clusters. Please log in to the Couchbase Capella UI where you'll be able to manage buckets and database users once your hosted cluster has been deployed.

~> **WARNING:** Updating cluster servers will cause your cluster to redeploy. You won't be able to access the cluster in the Couchbase Capella UI until it has been deployed. Your cluster will remain functional during this time.

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
    type     = "Basic"
  }
  servers {
    size     = 3
    compute  = "m5.xlarge"
    services = ["data"]
    storage {
      type = "GP3"
      iops = "3000"
      size = "50"
    }
  }
}
```

## Argument Reference

- `name` - (Required) The name of the cluster you want to create.
- `project_id` - (Required) The id of the project where your cluster will be created. (Cannot be changed via this Provider after creation.)
- `description` - (Optional) A description for the cluster.

### Place

- `single_az` - (Required) A boolean to describe if there is only a single availability zone. (Cannot be changed via this Provider after creation.)

~> **WARNING:** `single_az` Has to be true if the you select the "Basic" support package.

#### Hosted

- `provider` - (Required) The name of the cloud provider you want your cluster to be hosted in. (Cannot be changed via this Provider after creation.)
- `region` - (Required) A valid region for the cloud provider that you want you cluster to be hosted in. (Cannot be changed via this Provider after creation.)
- `cidr` - (Required) The cidr block. (Cannot be changed via this Provider after creation.)

### Support Package

- `timezone` - (Required) The time zone that you would like to receive support from. `ET`, Eastern Time, `GMT`, Greenwich Mean Time , `IST`, India Standard Time, `PT`, Pacific Time, are the available time zones that you can specify.
- `type` - (Required) The support plan that you would like for your Capella cluster. `Basic`, `DeveloperPro`, `Enterprise` are the available support plans that you can specify.

For more detailed information on support packages, you can view this [detailed plan comparison](https://www.couchbase.com/support-policy/cloud).

### Servers

- `size` - (Required) The number of nodes in your cluster.
- `compute` - (Required) The name of the compute instance type.
- `services` - (Required) A list of Couchbase services that you want in your cluster. `Data`, `Query`, `Index`, `Search`,`eventing`, `analytics` are the available services that you can specify.

#### Storage

- `type` - (Required) The name of the storage type. `GP3`, `IO2` are the available storage types that you can specify.
- `iops` - (Required) The number of the IOPS.

~> **IMPORTANT:** The minimum value is 3000 for GP3, and 1000 if IO2. The maximum value is 16000 for GP3, and 64000 if IO2.

- `size` - (Required) The storage per node in gigabytes.

~> **IMPORTANT:** The minimum storage per node is 50Gb. The maximum storage per node is 16Tb.

## Attribute Reference

- `id` - The cluster id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clustersv3).
