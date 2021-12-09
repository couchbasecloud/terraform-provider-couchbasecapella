resource "couchbasecapella_hosted_cluster" "cluster" {
  name        = var.hosted_cluster_name
  project_id  = couchbasecapella_project.project.id
  description = "Example Description"
  place {
    single_az = true
    hosted {
      provider = "aws"
      region   = "us-west-2"
      cidr     = "10.0.16.0/20"
    }
  }
  support_package {
    timezone = "GMT"
    type     = "Basic"
  }
  servers {
    size     = 3
    compute  = "m5.xlarge"
    services = ["data"]
    storage {
      type = "GP3"
      iops = "3000"
      size = "50"
    }
  }
}
