---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "luminate_collection_site_link Resource - terraform-provider-luminate"
subcategory: ""
description: |-
  
---

# luminate_collection_site_link (Resource)



## Example Usage

```terraform
# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_collection_site_link" "new-collection-site-link" {
  site_id        = "c11e4576-53c8-4617-a408-5d31a9c9e954"
  collection_ids = sort(["8d945145-0d0a-4b76-b6a7-8f7af4fc8dc3"])
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `collection_ids` (List of String) Collection IDs
- `site_id` (String) Site ID

### Read-Only

- `id` (String) The ID of this resource.
