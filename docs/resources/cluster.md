# couchbasecapella_cluster Resource

`couchbasecapella_cluster` provides a Cluster resource. The resource lets you create, edit and delete clusters. The resource requires your Project ID.

## Example Usage

### Example AWS cluster

```hcl
resource "couchbasecapella_cluster" "test" {
  name       = "cluster_name"
  cloud_id   = "your_cloud_id"
  project_id = "your_project_id"
  servers {
    size     = 3
    services = ["data", "query", "index"]
    aws {
      instance_size = "m5.xlarge"
      ebs_size_gib  = 50
    }
  }
}
```

### Example Azure cluster.

```hcl
resource "couchbasecapella_cluster" "test" {
  name       = "cluster_name"
  cloud_id   = "your_cloud_id"
  project_id = "your_project_id"
  servers {
    size     = 3
    services = ["data", "query", "index"]
    azure {
      instance_size = "Standard_F4s_v2"
      volume_type  = "P4"
    }
  }
}
```

## Argument Reference

- `name` - (Required) The name of the cluster you want to create. (Cannot be changed via this Provider after creation.)
- `cloud_id` - (Required) The id of the cloud where your cluster will be created. (Cannot be changed via this Provider after creation.)
- `project_id` - (Required) The id of the project where your cluster will be created. (Cannot be changed via this Provider after creation.)

### Servers

- `size` - (Required) The number of nodes in your cluster.
- `services` - (Required) A list of Couchbase services that you want in your cluster. (Cannot be changed via this Provider after creation.)

#### AWS

- `instance_size` - (Required) The name of the aws instance type. (Cannot be changed via this Provider after creation.)
- `ebs_size_gib` - (Required) The size of the ebs volume in gigabytes. (Cannot be changed via this Provider after creation.)

#### Azure

- `instance_size` - (Required) The name of the azure instance type. (Cannot be changed via this Provider after creation.)
- `volume_type` - (Required) The name of the azure volume type. (Cannot be changed via this Provider after creation.)

## Attribute Reference

- `id` - The cluster id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
