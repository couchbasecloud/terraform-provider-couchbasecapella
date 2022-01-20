resource "couchbasecapella_project" "project" {
  name = var.project_name
}
output "terraform_project_id" {
  value = couchbasecapella_project.project.id
}
