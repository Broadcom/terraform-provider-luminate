# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

data "luminate_group" "my-groups" {
  identity_provider_id = "identity_provider_id"
  groups               = ["group1", "group2"]
}
