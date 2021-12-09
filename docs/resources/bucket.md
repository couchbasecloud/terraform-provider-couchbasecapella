# couchbasecapella_bucket Resource

`couchbasecapella_buckets` provides a Bucket resource. The resource allows buckets to be created, edited and deleted inside a cluster. This resource requires your VPC Cluster ID.

~> **WARNING:** This current release of Terraform Couchbase Capella Provider doesn't support updating the bucket. Please log in to the Couchbase Capella UI where you'll be able to edit the memory quota and bucket access for each bucket.

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
