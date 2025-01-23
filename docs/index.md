Basic configuration and usage
==========

Symantec ZTNA terraform provider is used to create and
manage resources supported by Symantec ZTNA platform.

Provider configuration
-----------

To use the provider it must first be configured to access Symantec ZTNA management API.

#### Example Usage

```terraform
# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

provider "luminate" {
  api_endpoint = "api.example.luminatesec.com"
}
```

API Endpoint
------

The API endpoint address is based on the tenant name in Symantec ZTNA

The format is as follows:
```
api.<tenant_name>.luminatesec.com
```
For example:
If the tenant name is "mycompany" the API endpoint address would be
"api.mycompany.luminatesec.com"

Authentication
-------

Authentication is done using an API Client credentials

#### Authenticate using environment variables

**shell**
```
$ export LUMINATE_API_CLIENT_ID=123456789
$ export LUMINATE_API_CLIENT_SECRET=abcdefghijk
```

**main.<span></span>tf**

```terraform
# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

provider "luminate" {
  api_endpoint = "api.example.luminatesec.com"
}
```

#### Authenticate using the provider block

```terraform
# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

provider "luminate" {
  api_endpoint      = "api.example.luminatesec.com"
  api_client_id     = "123456789"
  api_client_secret = "abcdefghijk"
}
```

  **Warning:** storing credentials in terraform files is not recommended and may lead to a secret leak in case the file is committed to a public repository

------

Provider usage example
-----------

This will create a site with one connector, web application and access
policy

```terraform
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
```
