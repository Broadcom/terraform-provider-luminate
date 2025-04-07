# **Terraform provider luminate**


#### Latest Binaries  

| Platform    |                                                                                                                                                                           |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Linux       | [terraform-provider-luminate-linux.zip](https://github.com/Broadcom/terraform-provider-luminate/releases/latest/download/terraform-provider-luminate-linux.zip)           |
| MacOS Intel | [terraform-provider-luminate-darwin.zip](https://github.com/Broadcom/terraform-provider-luminate/releases/latest/download/terraform-provider-luminate-darwin.zip)         |
| MacOS M1    | [terraform-provider-luminate-windows.zip](https://github.com/Broadcom/terraform-provider-luminate/releases/latest/download/terraform-provider-luminate- darwin_arm64.zip) |
| Windows     | [terraform-provider-luminate-windows.zip](https://github.com/Broadcom/terraform-provider-luminate/releases/latest/download/terraform-provider-luminate-windows.zip) <br/> |

[![CircleCI](https://circleci.com/gh/Broadcom/terraform-provider-luminate/tree/master.svg?style=shield)](https://circleci.com/gh/Broadcom/terraform-provider-luminate)  
---

#### Documentation

[Basic configuration and usage](#basic-configuration-and-usage)
- [Provider configuration](#provider-configuration)
- [API Endpoint](#api-endpoint)
- [Authentication](#authentication)
- [Usage Example](#provider-usage-example)

[Core resources](#core-resources)
- [Resource: luminate_site](#resource-luminate_site)
- [Resource: luminate_connector](#resource-luminate_connector)

[Application Resources](#application-resources)
- [Resource: luminate_web_application](#resource-luminate_web_application)
- [Resource: luminate_ssh_application](#resource-luminate_ssh_application)
- [Resource: luminate_rdp_application](#resource-luminate_rdp_application)
- [Resource: luminate_tcp_application](#resource-luminate_tcp_application)
- [Resource: luminate_ssh_gw_application](#resource-luminate_ssh_application)
- [Resource: luminate_segment_application](#resource-luminate_segment_application)

[Policy resources](#policy-resources)
- [Resource: luminate_rdp_access_policy](#resource-luminate_rdp_access_policy)
- [Resource: luminate_ssh_access_policy](#resource-luminate_ssh_access_policy)
- [Resource: luminate_web_access_policy](#resource-luminate_web_access_policy)
- [Resource: luminate_tcp_access_policy](#resource-luminate_tcp_access_policy)
- [Resource: luminate_web_activity_policy](#resource-luminate_web_activity_policy)

[Collection resources](#collection-resources)
- [Resource: luminate_collection](#resource-luminate_collection)
- [Resource: luminate_collection_site_link](#resource-luminate_collection_site_link)
- [Resource: luminate_tenant_role](#resource-luminate_tenant_role)
- [Resource: luminate_collection_role](#resource-luminate_collection_role)
- [Resource: luminate_site_role](#resource-luminate_site_role)

[Identities resources](#identities-resources)
- [Resource: luminate_group_user](#resource-luminate_group_user)

[Dns Resiliency resources](#DNS-Resiliency-resources)
- [Resource: luminate_dns_group_resiliency](#resource-luminate_dns_resiliency_group)
- [Resource: luminate_dns_server_resiliency](#resource-luminate_dns_resiliency_server)

[Data sources](#data-sources)
- [Data Source: luminate_identity_provider](#data-Source: luminate_identity_provider)
- [Data Source: luminate_user](#data-Source: luminate_user)
- [Data Source: luminate_group](#data-Source: luminate_group)
- [Data Source: luminate_aws_integration](#data-Source: luminate_aws_integration)
- [Data Source: luminate_ssh_client](#data-Source: luminate_ssh_client)


Basic configuration and usage
==========

Broadcom secure access cloud terraform provider is used to create and
manage resources supported by Secure access cloud platform.

Provider configuration
-----------

To use the provider it must first be configured to access Secure access
cloud management API.

#### Example Usage

```
provider "luminate" {         
    api_endpoint = "api.example.luminatesec.com"
}
```                            

API Endpoint
------

The API endpoint address is based on the tenant name in Secure access
cloud

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
```
provider "luminate" {|
    api_endpoint = "api.example.luminatesec.com"
}
```
#### Authenticate using the provider block

```
provider "luminate" {
  api_endpoint = "api.example.luminatesec.com"
  api_client_id = "123456789"
  api_client_secret = "abcdefghijk"
}
```
  **Warning:** storing credentials in terraform files is not recommended and may lead to a secret leak in case the file is committed to a public repository

------

Provider usage example
-----------

This will create a site with one connector, web application and access
policy

```
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
  name = "connector-${luminate_site.site.name}"
  site_id = "${luminate_site.site.id}"
  type = "linux"
}

#Create web application
resource "luminate_web_application" "nginx-app" {
  name = "nginx"
  site_id = "${luminate_site.site.id}"
  internal_address = "http://127.0.0.1:8080"
}

#Retrieve the id of local IDP
data "luminate_identity_provider" "idp" {
  identity_provider_name = "local"
}

#Retrieve users from IDP
data "luminate_user" "users" {
  identity_provider_id = "${data.luminate_identity_provider.idp.identity_provider_id}"
  users = ["local-user"]
}

#Create access policy and attach to application
resource "luminate_web_access_policy" "nginx-access-policy" {
  name = "nginx-access-policy"
  identity_provider_id = "${data.luminate_identity_provider.idp.identity_provider_id}"
  user_ids = [${data.luminate_user.users.user_ids}]
  applications = [${luminate_web_application.nginx-app.id}]
}

#One time command to start the created connector
output "run-command" {
  value = "${luminate_connector.connector.command}"
}
```


Migrating existing code to terraform 0.12
---------------------

Due to changes to HCL terraform code written for previous versions has
to be converted to new language version.

Terraform provides a built-in command to make the required changes.

For more information and detailed instructions refer to:
<https://www.terraform.io/upgrade-guides/0-12.html>

Usage
---
```
cd /terraform-repo
terraform 0.12upgrade
```

**NOTE:** 0.12upgrade sub command will change code in-place overwriting existing files.

# Core resources

Re­­­source: luminate_site
----------

Provides secure access cloud site resource

­­­

#### Example Usage

```
resource "luminate_site" "new-site" {
  name = "my-new-site"
}
```
#### Argument Reference

The following arguments are supported:

- **name** (Required) The name of the site

- **region** (Optional) Connectivity region. If not specified, the default region will be used

- **mute_health_notification** (Optional) Don't send notification
    if the site is down

- **kubernetes_persistent_volume_name** (Optional) Kubernetes
    persistent volume name - only relevant if running on top kubernetes

- **kerberos** - (Optional)

    -   **domain** (Required) - Active Directory domain name you want
        to SSO with.

    -   **kdc_address** - (Required) - The hostname of the primary
        domain controller.

    -   **keytab_pair** - (Required) - The absolute path of the keytab
        file

  
  **NOTE:** kerberos block is optional, but if specified all nested fields are required

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the site

#### Import
```
$ terraform import luminate_site.new-site site-id
```

Re­­­source: luminate_connector
------------

Provides secure access cloud connector resource

­­­

#### Example Usage
```
resource "luminate_connector" "connector" {
  name = "connector-name"
  site_id = "site-id"
  type = "linux"
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the connector

-   **site_id -** (Required) site id to attach the connector

-   **type -** (Required) type of the connector. Valid types: **linux**
    \| **kubernetes** \| **windows** \| **docker-compose**

**NOTE:** Connector resource is immutable. Changing any of the arguments will trigger recreation


#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the connector

-   **command** - command for deploying Luminate connector

-   **otp -** one time password for running Luminate connector

#### Import

```
$ terraform import luminate_connector.connector connector_id
```

Application Resources
==========

Resource: luminate_web_application
----------

Provides Secure access cloud web application

#### Example Usage

```
resource "luminate_web_application" "new-web-application" {
  name = "web-application"
  site_id = "site_id"
  internal_address = "http://127.0.0.1:8080"
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the applications

-   **site_id** - (Required) Site ID to which the application will be
    bound

-   **icon** - (Optional) Base64 representation of 40x40 icon

-   **visible** - (Optional) Indicates whether to show this application
    in the applications portal.

-   **notification_enabled** - (Optional) Indicates whether
    notifications are enabled for this application.

-   **subdomain** - (Optional) The application DNS subdomain.

-   **custom_external_address** - (Optional) The application custom
    DNS address that exposes the application.

-   **internal_address** - (Required) Internal address of the
    application, accessible by connector

-   **custom_root_path** - (Optional) Requests coming into the
    external address root path \'/\', will be redirected to this custom
    path instead.

-   **health_url** - (Optional) Health check path. The URI is relative
    to the external address.

-   **health_method** - (Optional) HTTP method to validate application
    health. Valid methods: GET \| HEAD

-   **default_content_rewrite_rules_enabled** - (Optional)
    Indicates whether to enable automatic translation of all occurrences
    of the application internal address to its external address on most
    prominent content types and relevant headers.

-   **default_header_rewrite_rules_enabled** - (Optional) Indicates
    whether to enable automatic translation of all occurrences of the
    application internal address to its external address on relevant
    headers.

-   **use_external_address_for_host_and_sni** - (Optional)
    Indicates whether to use external address for host header and SNI.

-   **linked_applications** - (Optional) This property should be set
    in a scenario where the defined application contains resources that
    reference additional web applications by their internal domain name.

-   **header_customization** - (Optional) Custom headers key:value
    pairs to be added to all requests.
- 
-   **collection_id -** (Optional) Collection id to be linked to app, if empty will be assigned to default collection


#### Attribute Reference


In addition to arguments above, the following attributes are exported:

-   **id** - id of the application

-   **external_address**

-   **luminate_address**

#### Import

```
$ terraform import luminate_web_application.new-web-application application_id
```

Re­­­source: luminate_ssh_application
-------

Provides Secure access cloud SSH application

­­­

#### Example Usage


```
resource "luminate_ssh_application" "new-ssh-application" {  
    site_id = "site_id"
    name = "ssh-applications"
    internal_address = "tcp://127.0.0.1:22"
}                                                   
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the applications

-   **site_id** - (Required) Site ID to which the application will be
    bound

-   **icon** - (Optional) Base64 representation of 40x40 icon

-   **visible** - (Optional) Indicates whether to show this application
    in the applications portal.

-   **notification_enabled** - (Optional) Indicates whether
    notifications are enabled for this application.

-   **subdomain** - (Optional) The application DNS subdomain.

-   **custom_external_address** - (Optional) The application custom
    DNS address that exposes the application.

-   **internal_address** - (Required) Internal address of the
    application, accessible by connector

-   **collection_id -** (Optional) Collection id to be linked to app, if empty will be assigned to default collection

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the application

-   **external_address**

-   **luminate_address**

#### Import

```
$ terraform import luminate_ssh_application.new-ssh-application  application_id
```


Re­­­source: luminate_rdp_application
------

Provides Secure access cloud RDP application

­­
#### Example Usage

```
resource "luminate_rdp_application" "new-rdp-application" {
  site_id = "site_id"
  name = "rdp-application"
  internal_address = "tcp://127.0.0.1"
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the applications

-   **site_id** - (Required) Site ID to which the application will be
    bound

-   **icon** - (Optional) Base64 representation of 40x40 icon

-   **visible** - (Optional) Indicates whether to show this application
    in the applications portal.

-   **notification_enabled** - (Optional) Indicates whether
    notifications are enabled for this application.

-   **subdomain** - (Optional) The application DNS subdomain.

-   **custom_external_address** - (Optional) The application custom
    DNS address that exposes the application.

-   **internal_address** - (Required) Internal address of the
    application, accessible by connector

- **collection_id -** (Optional) Collection id to be linked to app, if empty will be assigned to default collection

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the application

-   **external_address**

-   **luminate_address**

#### Import

```
$ terraform import luminate_rdp_application.new-rdp-application application_id
```

Re­­­source: luminate_tcp_application
-----------

Provides Secure access cloud TCP application

­­­

#### Example Usage

```
resource "luminate_tcp_application" "new-tcp-application" {
  name = "tcp-application"
  site_id = "site-id"
  target {
    address = "127.0.0.1"
    ports = ["8080"]
  }
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the applications

-   **site_id** - (Required) Site ID to which the application will be
    bound

-   **icon** - (Optional) Base64 representation of 40x40 icon

-   **visible** - (Optional) Indicates whether to show this application
    in the applications portal.

-   **notification_enabled** - (Optional) Indicates whether
    notifications are enabled for this application.

-   **subdomain** - (Optional) The application DNS subdomain.

-   **custom_external_address** - (Optional) The application custom
    DNS address that exposes the application.

-   **target** - (Required) - list of TCP application targets

    -   **address** - (Required) application target address.

    -   **ports** - (Required) list of forwarded ports.
    
    -   **collection_id -** (Optional) Collection id to be linked to app, if empty will be assigned to default collection

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the application

-   external_address

-   luminate_address


#### Import

```
$ terraform import luminate_tcp_application.new-tcp-application  application_id
```

Re­­­source: luminate_ssh_gw_application
------------

Provides Secure access cloud SSH GW application

­­­

#### Example Usage

```
resource "luminate_ssh_gw_application" "new-sshgw-access" {  
  site_id = "site_id"
  name = "sshgw-application"

  integration_id = "integration_id",

  tags {
    Type = "ssh-gw-demo"
  }

  vpc {
    region = "eu-west-1"
    cidr_block = "172.31.0.0/16"
    vpc_id = "vpc-123456789"
  }
}
```
#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the applications

-   **site_id** - (Required) Site ID to which the application will be
    bound

-   **icon** - (Optional) Base64 representation of 40x40 icon

-   **visible** - (Optional) Indicates whether to show this application
    in the applications portal.

-   **notification_enabled** - (Optional) Indicates whether
    notifications are enabled for this application.

-   **subdomain** - (Optional) The application DNS subdomain.

-   **custom_external_address** - (Optional) The application custom
    DNS address that exposes the application.

-   **internal_address** - (Required) Internal address of the
    application, accessible by connector

-   **integration_id** - (Required) integration id used to set up the
    ssh gw application

-   **tags** - (Required) a map of tags used to determine which
    machines is included as part of this ssh gw

-   **vpc** - (Required) A list of vpc definitions used to determine
    the target group to include as part of the ssh gw application

    -   **vpc_id** - (Required) - the vpc id of the vpc containing
        target machines

    -   **region** - (Required) - the region containing the target
        machines

    -   **cidr_block** - (Required) - the cidr block of the machines
        to include

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the application

-   **segment_id**

-   **external_address**

-   **luminate_address**

#### Import

```
$ terraform import luminate_ssh_gw_application.new-sshgw-access  application_id
```

Resource: luminate_segment_application
------------

Provides Secure access cloud Segment application


#### Example Usage

```
resource "luminate_segment_application" "nginx-app" {
  name = "nginx"
  site_id = "${luminate_site.site.id}"
  segment_settings {
	original_ip = "10.60.30.0/24"
	}
}
```
Or using multiple_segment_settings
```
resource "luminate_segment_application" "new-segment-application" {
	name = "ngnix"
	sub_type = "SEGMENT_SPECIFIC_IPS"
	site_id = "${luminate_site.new-site.id}"
  	multiple_segment_settings {
            original_ip = ["192.168.1.1", "192.168.1.2"]
  	}
}

```

#### Argument Reference

The following arguments are supported:

- **name -** (Required) name of the applications

- **site_id** - (Required) Site ID to which the application will be
    bound

- **segment settings** - The segment application settings. This field will be deprecated, please use multiple segment settings instead.

    - **original_ip** - (Required) The internal resource IP address which is used by the connector for access to the application.
  
- **multiple segment settings** - (Required) The segment application settings

    - **original_ip** - (Required) The internal resource IPs addresses which is used by the connector for access to the application.

Policy resources
============

Re­­­source: luminate_rdp_access_policy
---------------

Provides Secure access cloud RDP access policy

­­­

#### Example Usage

```
resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  name =  "my-rdp-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]

  validators = {
    web_verification = true
  }

  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    location = ["Wallis and Futuna"]
  }
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the policy

-   **enabled -** (Optional) Indicates whether this policy is enabled.

-   **identity_provider_id -** (Optional) The identity provider id

-   **user_ids -** (Optional) The user entities to which this policy
    applies.

-   **group_ids -** (Optional) The group entities to which this policy
    applies.

-   **applications** - (Required) The applications to which this policy
    applies.

-   **validators** - (Optional)

    -   **web_verification** - (Optional) Indicate whatever to perform
        web verification validation. not compatible for HTTP
        applications

-   **conditions** - (Optional)

    -   **location** - (Optional) - location based condition, specify
        the list of accepted locations.

    -   **source_ip** - (Optional) - source ip based condition, specify
        the allowed CIDR for this policy.

-   **allow_long_term_password** - (Optional) Indicates whether to
    allow long term password.

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the policy

#### Import

```
$ terraform import luminate_rdp_access_policy.new-rdp-access-policy  policy_id
```

Re­­­source: luminate_ssh_access_policy
------------

Provides Secure access cloud SSH access policy

­­­

#### Example Usage

```
resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
  name =  "my-ssh-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]
  accounts = ["ubuntu", "ec2-user"]
  allow_temporary_token = true
}
```
#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the policy

-   **enabled -** (Optional) Indicates whether this policy is enabled.

-   **identity_provider_id -** (Optional) The identity provider id

-   **user_ids -** (Optional) The user entities to which this policy
    applies.

-   **group_ids -** (Optional) The group entities to which this policy
    applies.

-   **applications** - (Required) The applications to which this policy
    applies.

-   **accounts** - (Required) SSH/Unix accounts with which IDP entities
    and/or Luminate local users can access the SSH Server

-   **use_auto_mapping** - (Optional) Determine the strategy for
    mapping IDP entities to SSH/Unix accounts, and specifically indicate
    whether automatic mapping based on the logged-in IDP entity username
    is allowed. In case this property is set to TRUE, manually entered
    SSH accounts are ignored. This property is relevant for SSH
    applications only

-   **allow_agent_forwarding** - (Optional) Indicates whether SSH
    agent forwarding is allowed for transparent secure access to all
    corporate SSH Servers via this SSH application that acts a Bastion.
    This property is relevant for SSH applications only.

-   **allow_temporary_token** - (Optional) Indication whether
    authentication using a temporary token is allowed.

-   **allow_public_key** - (Optional) Indication whether
    authentication using long term secret is allowed.

-   **validators** - (Optional)

    -   **web_verification** - (Optional) Indicate whatever to perform
        web verification validation. not compatible for HTTP
        applications

-   **conditions** - (Optional)

    -   **location** - (Optional) - location based condition, specify
        the list of accepted locations.

    -   **source_ip** - (Optional) - source ip based condition, specify
        the allowed CIDR for this policy.

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the policy

#### Import

```
$ terraform import luminate_ssh_access_policy.new-ssh-access-policy  policy_id
```

Re­­­source: luminate_web_access_policy
---------

Provides Secure access cloud HTTP access policy

­­­

#### Example Usage

```
resource "luminate_web_access_policy" "new-web-access-policy" {
  name =  "my-web-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]
  
  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16", "8.8.8.8/24"]
    location = ["Wallis and Futuna"]

    managed_device = {
      symantec_cloudsoc = true
      symantec_web_security_service = false
    }
    validators {
        mfa = true
    }
  }
}
```
#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the policy

-   **enabled -** (Optional) Indicates whether this policy is enabled.

-   **identity_provider_id -** (Optional) The identity provider id

-   **user_ids -** (Optional) The user entities to which this policy
    applies.

-   **group_ids -** (Optional) The group entities to which this policy
    applies.

-   **applications** - (Required) The applications to which this policy
    applies.
- 
-   **validators** - (Optional)

    -   **mfa** - (Optional) Specifies whether to carry out mfa (multi-factor authentication) validation.

-   **conditions** - (Optional)

    -   **location** - (Optional) - location based condition, specify
        the list of accepted locations.

    -   **source_ip** - (Optional) - source ip based condition, specify
        the allowed CIDR for this policy.

    -   **managed_device** - (Optional) Indicate whatever to restrict
        access to managed devices only

        -   **opswat** - (Optional) Indicate whatever to restrict
            access to Opswat MetaAccess

        -   **symantec_cloudsoc** - (Optional) Indicate whatever to
            restrict access to symantec cloudsoc

        -   **symantec_web_security_service** - (Optional) Indicate
            whatever to restrict access to symantec web security service

    -   **unmanaged_device** - (Optional) Indicate whatever to
        restrict access to unmanaged devices only

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the policy

#### Import

```
$ terraform import luminate_web_access_policy.new-web-access-policy  policy_id
```

Re­­­source: luminate_tcp_access_policy
---------

Provides Secure access cloud TCP access policy

­­­

#### Example Usage

```
resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  name =  "my-tcp-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]
  accounts = ["ubuntu", "ec2-user"]
  allow_temporary_token = true
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the policy

-   **enabled -** (Optional) Indicates whether this policy is enabled.

-   **identity_provider_id -** (Optional) The identity provider id

-   **user_ids -** (Optional) The user entities to which this policy
    applies.

-   **group_ids -** (Optional) The group entities to which this policy
    applies.

-   **applications** - (Required) The applications to which this policy
    applies.

-   **allow_temporary_token** - (Optional) Indication whether
    authentication using a temporary token is allowed.

-   **allow_public_key** - (Optional) Indication whether
    authentication using long term secret is allowed.

-   **validators** - (Optional)

    -   **web_verification** - (Optional) Indicate whatever to perform
        web verification validation. not compatible for HTTP
        applications

-   **conditions** - (Optional)

    -   **location** - (Optional) - location based condition, specify
        the list of accepted locations.

    -   **source_ip** - (Optional) - source ip based condition, specify
        the allowed CIDR for this policy.

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the policy

#### Import

```
$ terraform import luminate_tcp_access_policy.new-tcp-access-policy  policy_id
```

Re­­­source: luminate_web_activity_policy
---------

Provides Secure access cloud HTTP activity policy

­­­

#### Example Usage

```
resource "luminate_web_activity_policy" "new-web-activity-policy" {
  name =  "my-web-activity-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]
  
  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16", "8.8.8.8/24"]
    location = ["Wallis and Futuna"]

    managed_device = {
      symantec_web_security_service = false
    }
  }
  
  rules = [
            {
              action = "BLOCK_USER"
              conditions = {
                uri_accessed = true
                http_command = true
                arguments = {
                  uri_list = ["/admin", "/users"]
                  commands = ["GET", "POST"]
                }
              }
            },
            {
              action = "DISCONNECT_USER"
              conditions = {
                file_uploaded = true
                file_downloaded = true
              }
            }
         ]
}
```
#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the policy

-   **enabled -** (Optional) Indicates whether this policy is enabled.

-   **identity_provider_id -** (Optional) The identity provider id

-   **user_ids -** (Optional) The user entities to which this policy
    applies.

-   **group_ids -** (Optional) The group entities to which this policy
    applies.

-   **applications** - (Required) The applications to which this policy
    applies.

-   **rules** - (Required) The constraints on the actions to perform
    upon user web activity (non-empty array of nested rule objects)

    -   **rule** - Activity rule object

        -   **action** (Required) - The action to apply, allowed values: 
            "BLOCK", "BLOCK_USER", "DISCONNECT_USER"
        
        -   **conditions** (Required) - The conditions to apply the action

            -   **file_downloaded** (Optional) Indicate whether File 
                Downloaded condition is enabled
            
            -   **file_uploaded** (Optional) Indicate whether File
                Uploaded condition is enabled

            -   **uri_accessed** (Optional) Indicate whether URI Access
                condition is enabled, requires the URI List argument

            -   **http_command** (Optional) Indicate whether HTTP Command
                condition is enabled, requires the Commands argument

            -   **arguments** (Optional) - The arguments for the enabled
                conditions, required only if related conditions are enabled

                -   **uri_list** (Optional) - The URI List argument, 
                    required for the URI Accessed condition if enabled
                
                -   **commands** (Optional) - The Commands argument, 
                    required for the HTTP Command condition if enabled
                

-   **conditions** - (Optional)

    -   **location** - (Optional) - location based condition, specify
        the list of accepted locations.

    -   **source_ip** - (Optional) - source ip based condition, specify
        the allowed CIDR for this policy.

    -   **managed_device** - (Optional) Indicate whatever to restrict
        policy to managed devices only

        -   **opswat** - (Optional) Indicate whatever to restrict
            policy to Opswat MetaAccess

        -   **symantec_cloudsoc** - (Optional) Indicate whatever to
            restrict policy to symantec cloudsoc

        -   **symantec_web_security_service** - (Optional) Indicate
            whatever to restrict policy to symantec web security service

    -   **unmanaged_device** - (Optional) Indicate whatever to
        restrict policy to unmanaged devices only

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **id** - id of the policy

#### Import

```
$ terraform import luminate_web_activity_policy.new-web-activity-policy  policy_id
```

Collection resources
============

Resource: luminate_collection
----------

Provides Secure access cloud collection resource

#### Example Usage

```
resource "luminate_collection" "new-collection" {
  name = "my-collection"
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the collection

Resource: luminate_collection_site_link
---------------

Provides Secure access cloud link between site and collection

#### Example Usage

```
resource "luminate_collection_site_link" "new-collection-site-link" {
      site_id = "c11e4576-53c8-4617-a408-5d31a9c9e954"
	  collection_ids = sort(["8d945145-0d0a-4b76-b6a7-8f7af4fc8dc3"])
	}
```

#### Argument Reference

The following arguments are supported:
-   **site_id -** (Required) Site id
-   **collection_ids -** (Required) Collection ids to be linked to site must be sorted



Resource: luminate_tenant_role
---------------

Provides Secure access cloud assign entity to tenant role

#### Example Usage

```
	resource "luminate_tenant_role" "tenant-admin" {
		role_type = "TenantAdmin"
		identity_provider_id =  "local"
		entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
		entity_type = "User"
	}
```

#### Argument Reference

The following arguments are supported:
-   **role_type  -** (Required) the role to assign TenantAdmin | TenantViewer
-   **identity_provider_id -** (Required) The identity provider id
-   **entity_id -** (Required) The entity id in idp
-   **entity_type -** (Required) the entity type User | Group | ApiClient

Resource: luminate_collection_role
---------------

Provides Secure access cloud assign entity to collection role

#### Example Usage

```
    resource "luminate_collection_role" "policy-owner" {
		role_type = "PolicyOwner"
		identity_provider_id =  "local"
		entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
		entity_type = "User"
		collection_id = "${luminate_collection.collection.id}"
	}   
```

#### Argument Reference

The following arguments are supported:
-   **role_type  -** (Required) the role to assign PolicyOwner | ApplicationOwner
-   **identity_provider_id -** (Required) The identity provider id
-   **entity_id -** (Required) The entity id in idp
-   **entity_type -** (Required) the entity type User | Group | ApiClient
-   **collection_id -** (Required) Collection id to be assigned

Resource: luminate_site_role
---------------

Provides Secure access cloud assign entity to site role

#### Example Usage

```
	resource "luminate_site_role" "site-editor" {
		role_type = "SiteEditor"
		identity_provider_id =  "local"
		entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
```

#### Argument Reference

The following arguments are supported:
-   **role_type  -** (Required) the role to assign SiteEditor | SiteConnectorDeployer
-   **identity_provider_id -** (Required) The identity provider id
-   **entity_id -** (Required) The entity id in idp
-   **entity_type -** (Required) the entity type User | Group | ApiClient
-   **site_id -** (Required) Site id to be assigned

# Identities resources

Resource: luminate_group_user
----------

Provides secure access cloud group_user resource

­­­

#### Example Usage

```
data "luminate_group" "my-groups" {
	identity_provider_id = "local"
	groups = ["group1"]
}

data "luminate_user" "my-users" {
	identity_provider_id = "local"
	users = ["user1"]
}

resource "luminate_group_user" "new_group_membership" {
	group_id = "${data.luminate_group.my-groups.group_ids.0}"
	user_id = "${data.luminate_user.my-users.user_ids.0}"
}
```
#### Argument Reference

The following arguments are supported:

-   **group_id -** (Required) Group id
-   **user_id -** (Required) User id to be assigned to group


# Integration resource

Resource: luminate_aws_integration
----------

Provides secure access cloud aws_integration resource

­
#### Example Usage

```
resource "luminate_aws_integration" "new-integration" {
	integration_name = "exampleIntegration"
}
```

#### Argument Reference

The following arguments are supported:

-   **integration_name -** (Required) name for the AWS integration

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **integration_id -** new integration id
-   **luminate_aws_account_id -** luminate AWS account ID
-   **aws_external_id -** the integration AWS external ID


# Integration Bind resource

Resource: luminate_aws_integration
----------

Provides secure access cloud aws_integration resource

­
#### Example Usage

```
resource "luminate_aws_integration" "new-integration" {
	integration_name = "exampleIntegrationBind"
}

//create and bind IAMrole and policy with new integration external ID and luminate account ID
resource "aws_iam_role" "test_role" {
  name = "exampleIntegrationBind"
  assume_role_policy = jsonencode({
	 Version= "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = "sts:AssumeRole"
                Condition = {
                    StringEquals = {
                        "sts:ExternalId": [
                            "${luminate_aws_integration.new-integration.aws_external_id}"
                        ]
                    }
                },
                Principal = {
                    "AWS" = [
                        "${luminate_aws_integration.new-integration.luminate_aws_account_id}"
                    ]
                }
            }
        ]
	})
}

resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
	  Sid = "VisualEditor0"
        Effect   = "Allow"
        Action = [
           "ec2:DescribeInstances",
           "ec2:DescribeVpcs",
           "ec2:DescribeRegions",
           "ec2:DescribeTags"
        ]
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role. test_role.name
  policy_arn = aws_iam_policy.policy.arn
}

resource "luminate_aws_integration_bind" "new-integration-bind" {
	integration_name = "${luminate_aws_integration.new-integration.integration_name}"
	integration_id= "${luminate_aws_integration.new-integration.integration_id}"
	aws_role_arn= "aws_iam_role_policy_attachment.test-attach.arn"
	luminate_aws_account_id= "${luminate_aws_integration.new-integration.luminate_aws_account_id}"	
	aws_external_id= "${luminate_aws_integration.new-integration.aws_external_id}"
	regions = ["us-west-1"]
}
```

#### Argument Reference

The following arguments are supported:

-   **integration_name -** (Required) name of the AWS integration
-   **integration_id -** (Required) ID of the AWS integration
-   **aws_role_arn -** (Required) AWS arn 
-   **luminate_aws_account_id -** (Required) luminate AWS account ID
-   **aws_external_id -** (Required) integration AWS external ID 
-   **regions -** (Required) regions to add

Re­­­source: luminate_DNS_server
------------

Provides secure access cloud DNS server

#### Example Usage
```
resource "luminate_site" "new-site" {
	name = "tfAccSite"
}
resource "luminate_dns_server" "new-dns" {
	site_id = "${luminate_site.new-site.id}"
	name = "testDNS"
	internal_address = "udp://10.0.0.1:53"
	dns_settings {
		domain_suffixes = ["company.com"]
	}
}
```

#### Argument Reference

The following arguments are supported:

-   **name -** (Required) name of the connector

-   **site_id -** (Required) site id to attach the connector

-   **internal_address** - (Required) Internal address of the
    application, accessible by connector
-   **dns_settings** [domain_suffixes] - (Required) The domain suffix


Data sources
==========

Data source: luminate_identity_provider
-----------

Use this resource to get an existing identity provider

­­­

#### Example Usage

```
data "luminate_identity_provider" "my-identity-provider" {
  identity_provider_name = "local"
}
```
#### Argument Reference

The following arguments are supported:

-   **identity_provider_name -** (Required) name of the identity
    provider

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **identity_provider_id** - id of the identity provider

Data Source: luminate_user
-------------

Use this resource to get one or more existing users

­­­

#### Example Usage

```
data "luminate_user" "my-users" {
  identity_provider_id = "identity_provider_id"
  users = ["user1@example.com", "user2@example.com"]
}
```
#### Argument Reference

The following arguments are supported:

-   **identity_provider_id -** (Required) id of the identity provider

-   **users -** (Required) List of usernames to retrieve

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **user_ids** - list of retrieved users ids

Data source: luminate_group
-----------

Use this resource to get one or more existing groups

­­­

#### Example Usage

```
data "luminate_group" "my-groups" {
  identity_provider_id = "identity_provider_id"
  groups = ["group1", "group2"]
}
```
#### Argument Reference

The following arguments are supported:

-   **identity_provider_id -** (Required) id of the identity provider

-   **groups -** (Required) List of group names to retrieve

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **group_ids** - list of retrieved groups ids

Data source: luminate_aws_integration
------------

Use this resource to retrieve an existing AWS integration

­­­

#### Example Usage

```
data "luminate_aws_integration" "my-integration" {
  integration_name = "integration_name"
}
```
#### Argument Reference

The following arguments are supported:

-   **integration_name -** (Required) name of an existing AWS
    integration

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

-   **integration_id** - id of retrieved AWS integration


Data source: luminate_ssh_client
------------

Use this resource to retrieve an existing ssh-client

#### Example Usage

```
data "luminate_ssh_client" "my-ssh-client" {
  name = "test"
}
```

#### Argument Reference

- **name** (String) ssh-client to retrieve

#### Attribute Reference

In addition to arguments above, the following attributes are exported:

- **id** (String)
- **description** (String)
- **key_size** (Number)
- **expires** (String)
- **last_accessed** (String)
- **created_on** (String)
- **modified_on** (String)

# DNS Resiliency resources

Resource: luminate_dns_server_resiliency
----------

Provides CRUD of dns resiliency servers

­­­

#### Example Usage

```
resource "luminate_site" "new-site" {
	name = "tfAccSite"
}

resource "luminate_dns_group_resiliency" "new-dns-group" {
	name = "testDNSGroupResiliency"
	sendNotifications = true
	domainSuffixes = ["somedomain.com"]
}

data "luminate_dns_server_resiliency" "new-dns-server-resiliency" {
    name = "testDNSServerResiliency"
	site_id = "${luminate_site.new-site.id}"
	group_id = "${luminate_dns_group_resiliency.new-dns-group.id}"
	internal_address = "udp://20.0.0.1:63"
}

```
#### Argument Reference

The following arguments are supported:

-   **group_id -** (Required) Group id
-   **name -** (Required) Dns Server name
-   **site_id -** (Required) Associated Site id
-   **internal_address -** (Required) Dns server address 


Resource: luminate_dns_group_resiliency
----------

Provides crud of dns resiliency groups

­­­

#### Example Usage

```

resource "luminate_dns_group_resiliency" "new-dns-group" {
	name = "testDNSGroupResiliency"
	sendNotifications = true
	domainSuffixes = ["somedomain.com"]
}

```
#### Argument Reference

The following arguments are supported:

-   **name -** (Required) Dns Group name
-   **sendNotifications -** (Required) Indicates if notification are enabled
-   **domainSuffixes -** (Required) List of domain suffixes


#### Confluence page
https://fireglass.atlassian.net/wiki/x/dICL1