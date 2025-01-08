# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
  name = "my-ssh-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids             = ["user1_id", "user2_id"]
  group_ids            = ["group1_id", "group2_id"]

  applications          = ["application1_id", "application2_id"]
  accounts              = ["ubuntu", "ec2-user"]
  allow_temporary_token = true
}
