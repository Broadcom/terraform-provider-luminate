# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

data "luminate_user" "my-users" {
  identity_provider_id = "identity_provider_id"
  users                = ["user1@example.com", "user2@example.com"]
}
