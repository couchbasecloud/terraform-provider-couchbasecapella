terraform {
  required_providers {
    couchbasecloud = {
      source  = "terraform.couchbase.com/local/couchbasecloud"
      version = "1.0.0"
    }
  }
}

provider "couchbasecloud" {
  access_key  = var.access_key
  secret_key = var.secret_key
}

# resource "couchbasecloud_project" "test" {
#   name = "terraform_project"
# }


# resource "couchbasecloud_project" "test2" {
#   name = "terraform_project2"
# }

resource "couchbasecloud_cluster" "terraform_cluster" {
  name = "terraform_cluster2"
  cloud_id = var.cloud_id
  project_id = "949d63c9-490e-468d-af03-601cf632574f"
  servers {
      size= 3
      services = ["data"]
      aws {
          instance_size = "m5.xlarge"
          ebs_size_gib = 50
      }
    }
  
}
