variable "project_name" {
  type    = string
  default = "expense-system"
}

variable "domain_name" {
  type    = string
  default = "expense.mleone.dev"
}

variable "google_oauth_client_id" {
  type      = string
  sensitive = true
  default   = "fakeId"
}

variable "google_oauth_client_secret" {
  type      = string
  sensitive = true
  default   = "fakeSecret"
}
