# couchbasecapella_project Resource

`couchbasecapella_project` provides a Project resource. This resource allows projects to be created or deleted.

~> **WARNING:** Changing the name of an existing Project in your Terraform configuration will result the destruction of that Project the re-creation of the Project with the new name.

<!-- Projects that contain clusters cannot be destroyed without the associated clusters being destroyed first.  -->

Terraform will inform you of the destroyed/created resources before applying so be sure to verify any change to your environment before applying.

## Example Usage

```hcl
resource "couchbasecapella_project" "test" {
  name = "project_name"
}
```

## Argument Reference

- `name` - (Required) The name of the project you want to create. (Cannot be changed via this Provider after creation.)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The project id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#projects).
