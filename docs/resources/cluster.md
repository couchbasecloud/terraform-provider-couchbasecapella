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
- `services` - (Required) A list of Couchbase services that you want in your cluster. `Data`, `Query`, `Index`, `Search`,`eventing`, `analytics` are the available services that you can specify. (Cannot be changed via this Provider after creation.)

#### AWS

- `instance_size` - (Required) The name of the aws instance type. `m5.xlarge`, `m5.2xlarge`, `m5.4xlarge`, `m5.8xlarge`, `m5.12xlarge`, `m5.16xlarge`, `m5.24xlarge`, `r5.xlarge`, `r5.2xlarge`, `r5.4xlarge`, `r5.8xlarge`, `r5.12xlarge`, `r5.24xlarge` ,`c5.2xlarge`, `c5.4xlarge`, `c5.9xlarge`, `c5.12xlarge`, `c5.18xlarge`, `x1.16xlarge`, `x1.32xlarge` are the available AWS instance sizes that you can specify. (Cannot be changed via this Provider after creation.)
  For more information on AWS instance sizes, please visit the [AWS Documentation](https://aws.amazon.com/ec2/instance-types/)
- `ebs_size_gib` - (Required) The size of the ebs volume in gigabytes. (Cannot be changed via this Provider after creation.)

#### Azure

- `instance_size` - (Required) The name of the azure instance type.
  `Standard_F4s_v2`,`Standard_F8s_v2`,`Standard_F16s_v2`,`Standard_F32s_v2`,`Standard_F48s_v2`,`Standard_F64s_v2`,`Standard_F72s_v2`,`Standard_D4s_v3`, `Standard_D8s_v3`,`Standard_D16s_v3`,`Standard_D32s_v3`,`Standard_D48s_v3`,`Standard_D64s_v3`,`Standard_E4s_v3`,`Standard_E8s_v3`,`Standard_E16s_v3`,`Standard_E20s_v3`,`Standard_E32s_v3`,`Standard_E48s_v3`,`Standard_E64s_v3`,`Standard_GS2`,`Standard_GS3`,`Standard_GS4`,`Standard_GS5` are the available Azure instance sizes that you can specify. (Cannot be changed via this Provider after creation.)
  For more detailed information on azure instance types, please visit the [Azure Documentation](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).
- `volume_type` - (Required) The name of the azure volume type. `P4`, `P6`, `P10`, `P15`, `P20`, `P30`, `P40`, `P50`, `P60`, `P70` are the available volume types that you can specify. (Cannot be changed via this Provider after creation.)
  For more detailed information on volume types, please visit the [Azure Documentation](https://docs.microsoft.com/en-us/azure/virtual-machines/disks-types#premium-ssd-size).

## Attribute Reference

- `id` - The cluster id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
