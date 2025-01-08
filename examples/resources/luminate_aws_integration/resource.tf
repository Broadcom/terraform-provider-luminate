# Copyright (c) Symantec ZTNA
# SPDX-License-Identifier: MPL-2.0

resource "luminate_aws_integration" "new-integration" {
  integration_name = "exampleIntegrationBind"
}

//create and bind IAMrole and policy with new integration external ID and luminate account ID
resource "aws_iam_role" "test_role" {
  name = "exampleIntegrationBind"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Condition = {
          StringEquals = {
            "sts:ExternalId" : [
              "${luminate_aws_integration.new-integration.aws_external_id}"
            ]
          }
        },
        Principal = {
          "AWS" = [
            "${luminate_aws_integration.new-integration.luminate_aws_account_id}"
          ]
        }
      }
    ]
  })
}

resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "VisualEditor0"
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeVpcs",
          "ec2:DescribeRegions",
          "ec2:DescribeTags"
        ]
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.test_role.name
  policy_arn = aws_iam_policy.policy.arn
}

resource "luminate_aws_integration_bind" "new-integration-bind" {
  integration_name        = luminate_aws_integration.new-integration.integration_name
  integration_id          = luminate_aws_integration.new-integration.integration_id
  aws_role_arn            = "aws_iam_role_policy_attachment.test-attach.arn"
  luminate_aws_account_id = luminate_aws_integration.new-integration.luminate_aws_account_id
  aws_external_id         = luminate_aws_integration.new-integration.aws_external_id
  regions                 = ["us-west-1"]
}
