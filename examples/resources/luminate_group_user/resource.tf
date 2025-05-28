# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

data "luminate_group" "my-groups" {
  identity_provider_id = "local"
  groups               = ["group1"]
}

data "luminate_user" "my-users" {
  identity_provider_id = "local"
  users                = ["user1"]
}

resource "luminate_group_user" "new_group_membership" {
  group_id = data.luminate_group.my-groups.group_ids.0
  user_id  = data.luminate_user.my-users.user_ids.0
}
