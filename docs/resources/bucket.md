---
page_title: "Couchbase Capella: Bucket"
subcategory: ""
description: |-
Create, edit and delete Buckets in a Couchbase Capella Cluster.
---

# Resource couchbasecapella_bucket

`couchbasecapella_bucket` allows buckets to be created and deleted for a Couchbase Capella In-VPC Cluster. This resource requires the Cluster ID of an In-VPC Cluster.

~> **WARNING:** Changing the cluster_id, name or conflict_resolution of an existing Bucket in your Terraform configuration will result in the deletion and recreation of the Bucket with the new name in Capella. Before applying your changes, Terraform will inform you that it will destroy and recreate the resources. Make sure to review these changes before typing `yes` to apply them.

~> **VERY IMPORTANT:** **THIS MEANS YOU WILL LOSE ANY DATA IN THE EXISTING BUCKET**

## Example Usage

### Creating a Single Bucket

```hcl
resource "couchbasecapella_bucket" "test" {
  cluster_id          = "your_cluster_id"
  name                = "bucket_name"
  memory_quota        = "128"
  conflict_resolution = "seqno"
}
```

### Creating Multiple Buckets

Multiple instances of buckets should depend on each other using the field `depends_on`, as seen below. This tells Terraform to create buckets one after another, allowing enough time for the previous bucket creation job to be completed.

```hcl
resource "couchbasecapella_bucket" "test" {
  cluster_id          = "your_cluster_id"
  name                = "bucket_name"
  memory_quota        = "128"
  conflict_resolution = "seqno"
}

resource "couchbasecapella_bucket" "test2" {
  depends_on = [couchbasecapella_bucket.test]

  cluster_id          = "your_cluster_id"
  name                = "bucket_name_two"
  memory_quota        = "128"
  conflict_resolution = "seqno"
}
```

## Argument Reference

- `cluster_id` - (Required) The id of the cluster where your bucket will be created. This must be a valid UUID and an existing cluster id.
- `name` - (Required) The name of the bucket you want to create. The bucket name can contain letters, numbers, periods (.) or dashes (-). Bucket names cannot exceed 100 characters and must begin with a letter or a number.
- `memory_quota` - (Required) The amount of memory that the bucket will be allocated in megabytes. Buckets require a minimum of 100 MiB of memory per node.
- `conflict_resolution` - (Required) The type of conflict resolution. You can select `seqno`, sequence number, or `lww`, last write wins.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
