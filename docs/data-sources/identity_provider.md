---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "luminate_identity_provider Data Source - terraform-provider-luminate"
subcategory: ""
description: |-
  
---

# luminate_identity_provider (Data Source)



## Example Usage

```terraform
# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

data "luminate_identity_provider" "my-identity-provider" {
  identity_provider_name = "local"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `identity_provider_name` (String) The identity provider name as configured in Symantec ZTNA portal, if not specified local idp will be taken

### Read-Only

- `id` (String) The ID of this resource.
- `identity_provider_id` (String) A unique identifier of this Identity Provider
