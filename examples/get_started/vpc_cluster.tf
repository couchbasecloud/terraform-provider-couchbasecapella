resource "couchbasecapella_vpc_cluster" "cluster" {
  name       = var.vpc_cluster_name
  cloud_id   = var.cloud_id
  project_id = couchbasecapella_project.project.id
  servers {
    size     = 3
    services = ["data"]
    aws {
      instance_size = "m5.xlarge"
      ebs_size_gib  = 50
    }
  }
}
output "terraform_vpc_cluster_id" {
  value = couchbasecapella_vpc_cluster.cluster.id
}
