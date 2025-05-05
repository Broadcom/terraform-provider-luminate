# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_web_access_policy" "new-web-access-policy" {
  name = "my-web-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids             = ["user1_id", "user2_id"]
  group_ids            = ["group1_id", "group2_id"]

  applications = ["application1_id", "application2_id"]

  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16", "8.8.8.8/24"]
    location  = ["Wallis and Futuna"]

    managed_device = {
      symantec_cloudsoc             = true
      symantec_web_security_service = false
    }
    validators = {
      mfa = true
    }
  }
}
