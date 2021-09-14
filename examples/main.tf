terraform {
  required_providers {
    couchbasecloud = {
        source = "terraform.couchbase.com/local/couchbasecloud"
        version = "1.0.0"
    }
  }

 
}

provider "couchbasecloud" {
  acesss_key = "1jSX7XQRMN5qu4YZaB4SN5N4J0l497sz"
  secret_key = "647Q982Dv589XpMhWAafBM0EGcMWvN04TviBBIDzsPOlCo8qz8pTtlPHCcRkIuVd"
}

 resource "couchbase_cluster" "cluster_test" {
    name = "Cluster Test"
    cloudId = "8a65b442-edc4-458b-94de-999720376641"
    projectId = "3f9fff59-54f3-474a-b32e-b32fc71be81c"
}

