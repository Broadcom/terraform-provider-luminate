# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_site_role" "site-editor" {
  role_type            = "SiteEditor"
  identity_provider_id = "local"
  entity_id            = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
  entity_type          = "User"
  site_id              = luminate_site.site.id
}
