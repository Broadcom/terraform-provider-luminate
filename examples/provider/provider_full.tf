# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

#Configure the provider
provider "luminate" {
  api_endpoint = "api.example.luminatesec.com"
}

#Create site
resource "luminate_site" "site" {
  name = "my-new-site"
}

#Create connector and bind to site "my-new-site"
resource "luminate_connector" "connector" {
  name    = "connector-${luminate_site.site.name}"
  site_id = luminate_site.site.id
  type    = "linux"
}
#Create web application
resource "luminate_web_application" "nginx-app" {
  name    = "nginx"
  site_id = luminate_site.site.id
}

#Retrieve the id of local IDP
data "luminate_identity_provider" "idp" {
  identity_provider_name = "local"
}

#Retrieve users from IDP
data "luminate_user" "users" {
  identity_provider_id = data.luminate_identity_provider.idp.identity_provider_id
  users                = ["local-user"]
}

#Create access policy and attach to application
resource "luminate_web_access_policy" "nginx-access-policy" {
  name                 = "nginx-access-policy"
  identity_provider_id = data.luminate_identity_provider.idp.identity_provider_id
  user_ids             = ["${data.luminate_user.users.user_ids}"]
  applications         = ["${luminate_web_application.nginx-app.id}"]
}

#One time command to start the created connector
output "run-command" {
  value = luminate_connector.connector.command
}
