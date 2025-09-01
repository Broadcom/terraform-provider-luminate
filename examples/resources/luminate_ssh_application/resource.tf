# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_ssh_application" "new-ssh-application" {
  site_id          = "site_id"
  name             = "ssh-applications"
  internal_address = "tcp://127.0.0.1:22"
}