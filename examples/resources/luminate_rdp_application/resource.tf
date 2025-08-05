# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_rdp_application" "new-rdp-application" {
  site_id          = "site_id"
  name             = "rdp-application"
  internal_address = "tcp://127.0.0.1"
}
