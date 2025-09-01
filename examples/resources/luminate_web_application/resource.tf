# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_web_application" "new-web-application" {
  name             = "web-application"
  site_id          = "site_id"
  internal_address = "http://127.0.0.1:8080"
}