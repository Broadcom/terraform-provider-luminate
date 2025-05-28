# Copyright (c) Broadcom Inc.
# SPDX-License-Identifier: MPL-2.0

resource "luminate_site_registration_key_version" "new_site_registration_key_version" {
}

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
  site_id                         = luminate_site.new-site.id
  version                         = luminate_site_registration_key_version.new_site_registration_key_version.version
  revoke_existing_key_immediately = true
}
