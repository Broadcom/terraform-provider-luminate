# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_segment_application" "nginx-app" {
  name    = "nginx"
  site_id = luminate_site.site.id
  segment_settings {
    original_ip = "10.60.30.0/24"
  }
}

#Or using multiple_segment_settings

resource "luminate_segment_application" "new-segment-application" {
  name     = "ngnix"
  sub_type = "SEGMENT_SPECIFIC_IPS"
  site_id  = luminate_site.new-site.id
  multiple_segment_settings {
    original_ip = ["192.168.1.1", "192.168.1.2"]
  }
}