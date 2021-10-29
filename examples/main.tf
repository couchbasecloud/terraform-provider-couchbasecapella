terraform {
  required_providers {
    couchbasecapella = {
      source  = "terraform.couchbase.com/local/couchbasecloud"
      version = "1.0.0"
    }
  }
}

provider "couchbasecapella" {
  access_key = var.access_key
  secret_key = var.secret_key
}

# resource "couchbasecapella_project" "test" {
#   name = "terraform_project"
# }


# resource "couchbasecapella_project" "test2" {
#   name = "terraform_project2"
# }

resource "couchbasecapella_cluster" "terraform_cluster" {
  name       = "this_capella_cluster_was_created_by_terraform"
  cloud_id   = var.cloud_id
  project_id = "949d63c9-490e-468d-af03-601cf632574f"
  servers {
    size     = 3
    services = ["data"]
    aws {
      instance_size = "m5.xlarge"
      ebs_size_gib  = 50
    }
  }
}
