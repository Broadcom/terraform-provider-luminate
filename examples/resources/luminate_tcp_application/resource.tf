# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_tcp_application" "new-tcp-application" {
  name    = "tcp-application"
  site_id = "site-id"
  target {
    address = "127.0.0.1"
    ports   = ["8080"]
  }
}
