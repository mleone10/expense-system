# Terraform backend configuration
terraform {
  backend "s3" {
    bucket = "leone-terraform-states"
    key    = "expense-system.tfstate"
    region = "us-east-1"
  }
}

# AWS Provider configuration
provider "aws" {
  region = "us-east-1"
}

# AWS Cognito infrastructure
# Allows for simple authentication through third-party identity provider (Google)
resource "aws_cognito_user_pool" "users" {
  name = "${var.project_name}-users"
}

resource "aws_cognito_user_pool_domain" "domain" {
  domain          = "auth.${var.domain_name}"
  certificate_arn = aws_acm_certificate.client_certificate.arn
  user_pool_id    = aws_cognito_user_pool.users.id
}

resource "aws_route53_record" "hosted_zone_record" {
  name    = aws_cognito_user_pool_domain.domain.domain
  zone_id = aws_route53_zone.hosted_zone.zone_id
  type    = "A"

  alias {
    name                   = aws_cognito_user_pool_domain.domain.cloudfront_distribution_arn
    zone_id                = "Z2FDTNDATAQYW2" # This is the global CloudFront Distribution zone ID
    evaluate_target_health = false
  }
}

resource "aws_cognito_identity_provider" "google" {
  user_pool_id  = aws_cognito_user_pool.users.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    authorize_scopes              = "profile email openid"
    token_url                     = "https://www.googleapis.com/oauth2/v4/token"
    token_request_method          = "POST"
    oidc_issuer                   = "https://accounts.google.com"
    authorize_url                 = "https://accounts.google.com/o/oauth2/v2/auth"
    attributes_url                = "https://people.googleapis.com/v1/people/me?personFields="
    attributes_url_add_attributes = "true"
    client_id                     = var.google_oauth_client_id
    client_secret                 = var.google_oauth_client_secret
  }

  attribute_mapping = {
    email    = "email"
    username = "sub"
  }
}

resource "aws_cognito_user_pool_client" "client" {
  name                                 = "${var.project_name}-client"
  user_pool_id                         = aws_cognito_user_pool.users.id
  supported_identity_providers         = ["Google"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["email", "openid", "profile"]
  callback_urls                        = ["https://expense.mleone.dev"]
  generate_secret                      = true
  prevent_user_existence_errors        = "ENABLED"

  depends_on = [
    aws_cognito_identity_provider.google
  ]
}

# Static site infrastructure
# A few things happening here:
# * Create an S3 bucket that can only be accessed by CloudFront
# * Create an ACM certificate that uses DNS validation against a new Route53 hosted zone
# * Create a CloudFront distribution that can be accessed via our chosen domain
resource "aws_s3_bucket" "bucket" {
  bucket = "leone-${var.project_name}"
  acl    = "private"
}

resource "aws_cloudfront_origin_access_identity" "origin_access_identity" {
}

data "aws_iam_policy_document" "s3_iam_policy" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.bucket.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.origin_access_identity.iam_arn]
    }
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.bucket.id
  policy = data.aws_iam_policy_document.s3_iam_policy.json
}

resource "aws_acm_certificate" "client_certificate" {
  domain_name               = var.domain_name
  validation_method         = "DNS"
  subject_alternative_names = ["*.${var.domain_name}"]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_zone" "hosted_zone" {
  name = var.domain_name
}

resource "aws_route53_record" "validation_record" {
  for_each = {
    for data_validation_option in aws_acm_certificate.client_certificate.domain_validation_options : data_validation_option.domain_name => {
      name   = data_validation_option.resource_record_name
      record = data_validation_option.resource_record_value
      type   = data_validation_option.resource_record_type
    }
  }

  name            = each.value.name
  records         = [each.value.record]
  type            = each.value.type
  zone_id         = aws_route53_zone.hosted_zone.zone_id
  ttl             = 60
  allow_overwrite = true
}

resource "aws_acm_certificate_validation" "cert_validation" {
  certificate_arn         = aws_acm_certificate.client_certificate.arn
  validation_record_fqdns = [for record in aws_route53_record.validation_record : record.fqdn]
}

locals {
  s3_origin_id = "${var.project_name}-origin"
}

resource "aws_cloudfront_distribution" "cdn" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"
  aliases             = [var.domain_name]

  origin {
    domain_name = aws_s3_bucket.bucket.bucket_regional_domain_name
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.origin_access_identity.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = local.s3_origin_id
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.client_certificate.arn
    ssl_support_method  = "sni-only"
  }
}

resource "aws_route53_record" "cdn_record" {
  zone_id = aws_route53_zone.hosted_zone.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.cdn.domain_name
    zone_id                = aws_cloudfront_distribution.cdn.hosted_zone_id
    evaluate_target_health = false
  }
}

# AWS Lambda infrastructure
resource "aws_iam_role" "lambda_role" {
  name = "${var.project_name}-execution-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowLambdaToAssumeRole"
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_lambda_function" "lambda" {
  function_name = "${var.project_name}-api"
  role          = aws_iam_role.lambda_role.arn
  filename      = "handler.zip"
  handler       = "bin/lambdaserver"
  runtime       = "go1.x"
}
