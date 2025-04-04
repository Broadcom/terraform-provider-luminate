# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_tenant_role" "tenant-admin" {
  role_type            = "TenantAdmin"
  identity_provider_id = "local"
  entity_id            = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
  entity_type          = "User"
}
