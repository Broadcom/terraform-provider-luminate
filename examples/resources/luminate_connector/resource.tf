# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_connector" "connector" {
  name    = "connector-name"
  site_id = "site-id"
  type    = "linux"
}
