terraform {
  required_providers {
    couchbasecloud = {
      source  = "terraform.couchbase.com/local/couchbasecloud"
      version = "1.0.0"
    }
  }
}

provider "couchbasecloud" {
}

resource "couchbasecloud_project" "test" {
  name = "test_project1"
}


resource "couchbasecloud_project" "test2" {
  name = "test_project2"
}
