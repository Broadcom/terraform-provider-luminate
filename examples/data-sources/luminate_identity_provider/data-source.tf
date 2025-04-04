# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

data "luminate_identity_provider" "my-identity-provider" {
  identity_provider_name = "local"
}
