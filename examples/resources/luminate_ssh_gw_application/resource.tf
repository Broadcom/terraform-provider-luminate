# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_ssh_gw_application" "new-sshgw-access" {
  site_id        = "site_id"
  name           = "sshgw-application"
  integration_id = "integration_id"

  tags {
    Type = "ssh-gw-demo"
  }

  vpc {
    region     = "eu-west-1"
    cidr_block = "172.31.0.0/16"
    vpc_id     = "vpc-123456789"
  }
}
