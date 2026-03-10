resource "aws_acm_certificate" "viewer_certificate" {
  provider = aws.us_east_1
  domain_name = local.full_domain
  validation_method = "DNS"
  # Add subject alternative names if needed
  subject_alternative_names = [local.full_domain]

  lifecycle {
    create_before_destroy = true
  }
}

resource "terraform_data" "invalidation" {  # TODO: ensure this works
  triggers_replace = {
    distribution_id = aws_cloudfront_distribution.s3_distribution.id
    # Add files/paths that change to force invalidation
  }
  provisioner "local-exec" {
    command = "aws cloudfront create-invalidation --distribution-id ${aws_cloudfront_distribution.s3_distribution.id} --paths '/*'"
  }
}

resource "aws_cloudfront_origin_access_control" "cv-site-bucket" {
  name                              = "cv-origin-access-control"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  depends_on = [aws_s3_bucket.cv_site_bucket]
  origin {
    domain_name              = aws_s3_bucket.cv_site_bucket.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.cv-site-bucket.id
    origin_id                = local.s3_origin_id
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment = "Some comment" # TODO: change
  default_root_object = "index.html"
  price_class = "PriceClass_100" # TODO: ensure ok
  aliases = [local.full_domain]


  default_cache_behavior {
    allowed_methods = ["GET", "HEAD"]
    cached_methods = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id
    viewer_protocol_policy = "redirect-to-https"
    compress = true
    min_ttl = 0 # none
    default_ttl = 3600 # 1hr
    max_ttl          = 86400 # very long
    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.rewrite_urls.arn
    }
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  # # Cache behavior with precedence 0
  # ordered_cache_behavior {
  #   path_pattern           = "*.html"
  #   allowed_methods = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
  #   cached_methods = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
  #   target_origin_id       = local.s3_origin_id
  #   min_ttl                = 0
  #   default_ttl            = 86400
  #   max_ttl                = 31536000
  #   compress               = true
  #   viewer_protocol_policy = "redirect-to-https"
  #   forwarded_values {
  #     query_string = false # TODO: ok?
  #     cookies {
  #       forward = "none" # TODO: none?
  #     }
  #   }
  # }
  # # Cache behavior with precedence 1
  # ordered_cache_behavior {
  #   path_pattern           = "*.webp"
  #   allowed_methods = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
  #   cached_methods = ["GET", "HEAD"] # TODO: ["GET", "HEAD", "OPTIONS"]
  #   target_origin_id       = local.s3_origin_id
  #   min_ttl                = 0
  #   default_ttl            = 86400
  #   max_ttl                = 31536000
  #   compress               = true
  #   viewer_protocol_policy = "redirect-to-https"
  #   forwarded_values {
  #     query_string = false # TODO: ok?
  #     cookies {
  #       forward = "none" # TODO: none?
  #     }
  #   }
  # }

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
    acm_certificate_arn = aws_acm_certificate.viewer_certificate.arn # ARN of your ACM cert in us-east-1 # TODO: ???
    ssl_support_method = "sni-only"
  }
  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations = []
    }
  }
  tags = {
    component = "cv-site"
  }
}

resource "aws_cloudfront_function" "rewrite_urls" {
  name    = "rewrite-urls-function"
  runtime = "cloudfront-js-2.0"
  code    = file("${path.module}/scripts/cloudfront.js")
  # Automatically publish the function to the LIVE stage
  publish = true

  # Learn more about the aws_cloudfront_function resource in the [Terraform Registry](https://registry.terraform.io)
}