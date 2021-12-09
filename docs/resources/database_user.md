# couchbasecapella_database_user Resource

`couchbasecapella_database_user` provides a Database User resource. The resource allows Database Users to be created, edited and deleted inside a cluster. This resource requires your VPC Cluster ID.

## Example Usage

### All Bucket Access

```hcl
resource "couchbasecapella_database_user" "test" {
  cluster_id = "your_cluster_id"
  username = "username"
  password = "password"
  all_bucket_access = "data_reader"
}
```

### With Specific Bucket Access

```hcl
resource "couchbasecapella_database_user" "test" {
  cluster_id = "your_cluster_id"
  username = "username"
  password = "password"
  buckets{
    bucket_name = "your_bucket_name"
    bucket_access = ["data_reader", "data_writer"]
  }
}
```

## Argument Reference

- `cluster_id` - (Required) The id of the cluster where you will create your database user. (Cannot be changed via this Provider after creation.)
- `username` - (Required) The username of the database user you want to create. (Cannot be changed via this Provider after creation.)
- `password` - (Required) The password of the database user you want to create. Password must contain 8+ characters, 1+ upper case, 1+ numbers, 1+ symbols. (Cannot be changed via this Provider after creation.)

### Buckets

-> **WARNING:** You may only specify bucket level access for all buckets or specific buckets. Including both in your configuration will result in an error in creating your database user.

#### All Bucket Access

- `all_bucket_access` - (Required) The bucket level access you want the database user to have for all buckets. You can either specify `data_reader`, which will give read access, or `data_writer`, which will give read/write access. (Cannot be changed via this Provider after creation.)

#### Specific Bucket Access

- `bucket_name` - (Required) The name of the bucket that you want to specify access levels for. (Cannot be changed via this Provider after creation.)
- `bucket_access` - (Required) The bucket level access you want the database user to have for the named bucket. You can either specify `data_reader`, which will give read access, or `data_writer`, which will give read/write access. (Cannot be changed via this Provider after creation.)

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
