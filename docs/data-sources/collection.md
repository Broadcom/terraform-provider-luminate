---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "luminate_collection Data Source - terraform-provider-luminate"
subcategory: ""
description: |-
  Collection data source
---

# luminate_collection (Data Source)

Collection data source

## Example Usage

```terraform
# Copyright (c) Broadcom Inc.
# SPDX-License-Identifier: MPL-2.0

data "luminate_collection" "my-collection" {
  name = "my-collection-name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Collection name

### Read-Only

- `id` (String) Collection id
