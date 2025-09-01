# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_web_activity_policy" "new-web-activity-policy" {
  name = "my-web-activity-policy"

  identity_provider_id = "identity_provider_id"
  user_ids             = ["user1_id", "user2_id"]
  group_ids            = ["group1_id", "group2_id"]

  applications = ["application1_id", "application2_id"]

  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16", "8.8.8.8/24"]
    location  = ["Wallis and Futuna"]

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
        file_uploaded   = true
        file_downloaded = true
      }
    }
  ]
}