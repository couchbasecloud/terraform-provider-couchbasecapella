# couchbasecapella_allowlist Resource

`couchbasecapella_allowlist` provides an Allowlist resource. This resource allows IP address/range entries to be created, edited and deleted in the Allowlist for a cluster.

## Example Usage

### Example Temporary Entry to Allowlist

```hcl
resource "couchbasecapella_allowlist" "test" {
  cluster_id = "your_cluster_id"
  cidr_block = "cidr_block_to_add"
  rule_type  = "temporary"
  comment    = "comment"
  duration   = "2h0m0s"
}
```

### Example Permanent Entry to Allowlist

```hcl
resource "couchbasecapella_allowlist" "test" {
  cluster_id = "your_cluster_id"
  cidr_block = "cidr_block_to_add"
  rule_type  = "permanent"
  comment    = "comment"
}
```

## Argument Reference

- `cluster_id` - (Required) The id of the cluster where you will edit the allowlist. (Cannot be changed via this Provider after creation.)
- `cidr_block` - (Required) The IP address/range of addresses you want to add to the allowlist.
- `rule_type` - (Required) The rule type that you assign to the entry. You can either set this to `temporary` or `permanent` which will define how the entry exists in the allowlist.
- `comment` - (Required) A comment you would like to add to the entry.
- `duration` - (Required) The duration that the IP should exist in the allowlist. This will be ignored if you set the rule type to `permanent`.

<!-- ## Attribute Reference -->

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#clusters).
