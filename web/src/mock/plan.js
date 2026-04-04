/**
 * Mock for GET /api/v1/plan?namespace=...&name=...
 */
export const mockPlan = {
  planId: 'plan-20240115-103045',
  isDestroyPlan: false,
  isDriftDetectionPlan: false,
  content: `Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create
  ~ update in-place
  - destroy

Terraform will perform the following actions:

  # aws_s3_bucket.main will be updated in-place
  ~ resource "aws_s3_bucket" "main" {
        id     = "my-production-bucket"
      ~ tags   = {
          + "Environment" = "production"
            "Project"     = "tofu-controller"
        }
    }

  # aws_iam_role.lambda_exec will be created
  + resource "aws_iam_role" "lambda_exec" {
      + arn                   = (known after apply)
      + create_date           = (known after apply)
      + force_detach_policies = false
      + id                    = (known after apply)
      + managed_policy_arns   = (known after apply)
      + max_session_duration  = 3600
      + name                  = "lambda-execution-role"
      + path                  = "/"
      + unique_id             = (known after apply)
    }

Plan: 1 to add, 1 to change, 0 to destroy.`
}

export const mockDestroyPlan = {
  planId: 'plan-20240115-103045',
  isDestroyPlan: true,
  isDriftDetectionPlan: false,
  content: `Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # aws_vpc.main will be destroyed
  - resource "aws_vpc" "main" {
      - arn                              = "arn:aws:ec2:us-east-1::vpc/vpc-12345" -> null
      - cidr_block                       = "10.0.0.0/16" -> null
      - id                               = "vpc-12345" -> null
      - tags                             = {
          - "Name" = "production-vpc"
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.`
}

/**
 * Mock for GET /api/v1/plan?namespace=...&name=...
 * Returns destroy plan for resources named *-destroy or *-infra, normal plan otherwise.
 */
export function getMockPlan(endpoint) {
  const url = new URL(endpoint, 'http://localhost')
  const name = url.searchParams.get('name') || ''
  return name.includes('destroy') ? mockDestroyPlan : mockPlan
}
