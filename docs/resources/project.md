---
page_title: "Couchbase Capella: Project"
subcategory: ""
description: |-
Create, edit and delete Projects in Couchbase Capella.
---

# Resource couchbasecapella_project

`couchbasecapella_project` allows Projects to be created, edited and deleted in Couchbase Capella.

~> **WARNING:** Changing the name of an existing Project in your Terraform configuration will result in the deletion and recreation of the Project with the new name in Capella. Projects that contain clusters cannot be destroyed without the associated clusters being destroyed first. Before applying your changes, Terraform will inform you that it will destroy and recreate the resources. Make sure to review these changes before typing `yes` to apply them.

## Example Usage

```hcl
resource "couchbasecapella_project" "test" {
  name = "project_name"
}
```

## Argument Reference

- `name` - (Required) The name of the project you want to create.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The project id.

For more information see: [Couchbase Capella Public API Reference](https://docs.couchbase.com/cloud/reference/rest-endpoints-all.html#projects).
