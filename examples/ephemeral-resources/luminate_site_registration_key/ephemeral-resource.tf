# Copyright (c) Broadcom Inc.
# SPDX-License-Identifier: MPL-2.0

resource "luminate_site_registration_key_version" "new_site_registration_key_version" {
  version = 1
}

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
  is_applying                     = terraform.applying
  site_id                         = luminate_site.new-site.id
  revoke_existing_key_immediately = true
  should_rotate                   = luminate_site_registration_key_version.new_site_registration_key_version.version_changed
}
