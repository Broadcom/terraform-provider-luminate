# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_aws_integration_bind" "new-integration-bind" {
  integration_name        = luminate_aws_integration.new-integration.integration_name
  integration_id          = luminate_aws_integration.new-integration.integration_id
  aws_role_arn            = "aws_iam_role_policy_attachment.test-attach.arn"
  luminate_aws_account_id = luminate_aws_integration.new-integration.luminate_aws_account_id
  aws_external_id         = luminate_aws_integration.new-integration.aws_external_id
  regions                 = ["us-west-1"]
}
