resource "couchbasecapella_database_user" "database_user" {
  cluster_id        = couchbasecapella_cluster.cluster.id
  username          = var.dbuser
  password          = var.dbuser_password
  all_bucket_access = "data_writer"
}
