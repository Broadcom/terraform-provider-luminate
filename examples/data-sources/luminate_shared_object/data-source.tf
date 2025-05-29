# Copyright (c) Broadcom Inc.
# SPDX-License-Identifier: MPL-2.0

data "luminate_shared_object" "my-shared_object" {
  name = "my-shared-object"
  type = "IP_LIST"
}
