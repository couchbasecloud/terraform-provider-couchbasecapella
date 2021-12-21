---
page_title: "Couchbase Capella: Bucket"
subcategory: ""
description: |-
Create, edit and delete Buckets in a Couchbase Capella Cluster.
---

# Resource couchbasecapella_bucket

`couchbasecapella_bucket` allows buckets to be created, edited and deleted for a Couchbase Capella In-VPC Cluster. This resource requires the Cluster ID of an In-VPC Cluster.

~> **WARNING:** Changing the cluster_id, name or conflict_resolution of an existing Bucket in your Terraform configuration will result in the deletion and recreation of the Bucket with the new name in Capella. Before applying your changes, Terraform will inform you that it will destroy and recreate the resources. Make sure to review these changes before typing `yes` to apply them.

~> **VERY IMPORTANT:** **THIS MEANS YOU WILL LOSE ANY DATA IN THE EXISTING BUCKET**

## Example Usage

```hcl
resource "couchbasecapella_bucket" "bucket_test" {
  cluster_id          = "your_cluster_id"
  name                = "bucket_name"
  memory_quota        = "128"
  replicas            = "1"
  conflict_resolution = "seqno"
}
```

## Argument Reference

- `cluster_id` - (Required) The id of the cluster where your bucket will be created. (Cannot be changed via this Provider after creation.)
- `name` - (Required) The name of the bucket you want to create. (Cannot be changed via this Provider after creation.)
- `memory_quota` - (Required) The amount of memory that the bucket will be allocated in megabytes. Buckets require a minimum of 100 MiB of memory per node. (Cannot be changed via this Provider after creation.)
- `replicas` - (Required) The number of replicas this bucket will have. (Cannot be changed via this Provider after creation.)
- `conflict_resolution` - (Required) The type of conflict resolution. You can select `seqno`, sequence number, or `lww`, last write wins. (Cannot be changed via this Provider after creation.)

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
