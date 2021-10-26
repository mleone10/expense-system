terraform {
  backend "s3" {
    bucket = "leone-terraform-states"
    key    = "expense-system.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_cognito_user_pool" "users" {
  name = "${var.project_name}-users"
}

resource "aws_cognito_user_pool_domain" "domain" {
  domain       = "leoneexpensesystem"
  user_pool_id = aws_cognito_user_pool.users.id
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
    client_id                     = "${var.google_oauth_client_id}"
    client_secret                 = "${var.google_oauth_client_secret}"
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
  callback_urls                        = ["https://example.com/callback"]
  generate_secret                      = true
  prevent_user_existence_errors        = "ENABLED"

  depends_on = [
    aws_cognito_identity_provider.google
  ]
}

resource "aws_s3_bucket" "bucket" {
  bucket = "leone-${var.project_name}"
  acl    = "private"
}

locals {
  s3_origin_id = "${var.project_name}-origin"
}

resource "aws_cloudfront_distribution" "cdn" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"

  origin {
    domain_name = aws_s3_bucket.bucket.bucket_regional_domain_name
    origin_id   = local.s3_origin_id
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
    cloudfront_default_certificate = true
  }
}
