resource "aws_acm_certificate" "viewer_certificate" {
  domain_name       = var.domain # TODO: ensure ok and dont need the subdomain
  validation_method = "DNS"
  # Add subject alternative names if needed
  subject_alternative_names = ["*.${var.domain}",local.domain_name] # TODO: ensure ok
  # TODO: must be on US-EAST 1????

  lifecycle {
    create_before_destroy = true
  }
  options {
    export = "DISABLED" # Will incur costs otherwise
  }
}

resource "aws_cloudfront_origin_access_control" "cv-site-bucket" {
  name                              = "cv-origin-access-control"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}



resource "aws_cloudfront_origin_access_identity" "cv_site_bucket" {
  comment = "cv site bucket origin access identity"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name              = aws_s3_bucket.cv_site_bucket.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.cv-site-bucket.id
    origin_id                = local.s3_origin_id
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Some comment" # TODO: change
  default_root_object = "index.html"
  price_class = "PriceClass_100" # TODO: ensure ok
  aliases = [local.domain_name]
  # logging_config { # TODO: NO LOGS
  #   include_cookies = false
  #   bucket          = "mylogs.s3.amazonaws.com"
  #   prefix          = "myprefix"
  # }
  origin {
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.cv_site_bucket.cloudfront_access_identity_path
    }
    domain_name = local.domain_name # TODO: ok?
    origin_id   = local.s3_origin_id # TODO: ok?
  }

  default_cache_behavior {
    path_pattern = "*" # TODO: ok?
    allowed_methods  = ["GET", "HEAD"]# TODO: ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id
    viewer_protocol_policy = "redirect-to-https" # TODO: redirect-to-https?
    compress               = true # TODO: ok?
    min_ttl                = 0 # none
    default_ttl            = 3600 # 1hr
    max_ttl                = 86400 # very long
    headers      = ["Origin"] # TODO: ok?
    forwarded_values {
      query_string = false # TODO: ???
      headers      = ["Origin"] # TODO: ok?
      cookies {
        forward = "all" # TODO: ???
      }
    }
  }

  # Cache behavior with precedence 0
  ordered_cache_behavior {
    path_pattern     = "*.html"
    allowed_methods  = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
    target_origin_id = local.s3_origin_id
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = false # TODO: ok?
      cookies {
        forward = "none" # TODO: none?
      }
    }
  }
  # Cache behavior with precedence 1
  ordered_cache_behavior {
    path_pattern     = "*.webp"
    allowed_methods  = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
    target_origin_id = local.s3_origin_id
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = false # TODO: ok?
      cookies {
        forward = "none" # TODO: none?
      }
    }
  }

  # # Cache behavior with precedence 1
  # ordered_cache_behavior {
  #   path_pattern     = "/*"
  #   allowed_methods  = ["GET", "HEAD"]
  #   cached_methods   = ["GET", "HEAD"]
  #   target_origin_id = local.s3_origin_id
  #
  #   forwarded_values {
  #     query_string = true # TODO: ok?
  #
  #     cookies {
  #       forward = "all" # TODO: ok?
  #     }
  #   }
  #
  #   min_ttl                = 0
  #   default_ttl            = 3600
  #   max_ttl                = 86400
  #   compress               = true
  #   viewer_protocol_policy = "redirect-to-https"
  # }

  # viewer_certificate { # TODO: can we use this instead?
  #   cloudfront_default_certificate = true
  # }
  viewer_certificate {
    #cloudfront_default_certificate = false # TODO: ???
    acm_certificate_arn = aws_acm_certificate.viewer_certificate.arn # ARN of your ACM cert in us-east-1 # TODO: ???
    ssl_support_method  = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
  tags = {
    component = "cv-site"
  }
}


# Wait for the ACM certificate to be validated
resource "aws_acm_certificate_validation" "test" {
  certificate_arn         = aws_acm_certificate.viewer_certificate.arn
  # Use the FQDNs from the cloudflare_record resources to ensure Terraform waits for their creation
  validation_record_fqdns = [for record in cloudflare_dns_record.acm_validation : record.name] # TODO: is name ok instead of hostname?

  # TODO: Optional, Set a longer timeout if DNS propagation is slow
  timeouts {
    create = "10m"
  }
}