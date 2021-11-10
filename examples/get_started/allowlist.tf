resource "couchbasecapella_allowlist" "ip" {
  cluster_id = couchbasecapella_cluster.cluster.id
  cidr_block = var.ip_address
  rule_type  = "temporary"
  comment    = "Temporary IP address for accesing cluster"
  duration   = "2h0m0s"
}
