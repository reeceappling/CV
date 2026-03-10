resource "aws_acm_certificate" "viewer_certificate" {
  provider = aws.us_east_1
  domain_name = local.full_domain # TODO: ensure ok and dont need the subdomain
  validation_method = "DNS"
  # Add subject alternative names if needed
  subject_alternative_names = [local.full_domain] # TODO: ensure ok

  lifecycle {
    create_before_destroy = true
  }
  # options {
  #   export = "DISABLED" # Will incur costs otherwise # TODO: this
  # }
}

resource "terraform_data" "invalidation" {
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


# resource "aws_cloudfront_origin_access_identity" "cv_site_bucket" {
#   comment = "cv site bucket origin access identity"
# }

resource "aws_cloudfront_distribution" "s3_distribution" {
  depends_on = [aws_s3_bucket.cv_site_bucket] # TODO: ok?
  origin {
    domain_name              = aws_s3_bucket.cv_site_bucket.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.cv-site-bucket.id
    origin_id                = local.s3_origin_id # TODO: ok?
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment = "Some comment" # TODO: change
  default_root_object = "index.html"
  price_class = "PriceClass_100" # TODO: ensure ok
  aliases = [local.full_domain]
  # logging_config { # TODO: NO LOGS
  #   include_cookies = false
  #   bucket          = "mylogs.s3.amazonaws.com"
  #   prefix          = "myprefix"
  # }


  default_cache_behavior {
    allowed_methods = ["GET", "HEAD"]# TODO: ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id
    viewer_protocol_policy = "redirect-to-https" # TODO: redirect-to-https?
    compress = true # TODO: ok?
    min_ttl = 0 # none
    default_ttl = 3600 # 1hr
    max_ttl          = 86400 # very long
    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.rewrite_urls.arn
    }
    forwarded_values {
      query_string = false # TODO: ???
      # headers      = ["Origin"] # TODO: ok?
      cookies {
        forward = "none" # TODO: all???
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

locals {
  # The JavaScript code can be inline or loaded from a file
  rewrite_function_code = <<-EOT
    function handler(event) {
      var request = event.request;
      var uri = request.uri;

      if (path.length > 1){
        const lastSegment = path.split('/').filter(Boolean).pop() || '';
        if (lastSegment.length !==0 && !lastSegment.includes('.')){
          request.uri += '.html';
        }
      }

      return request;
    }
  EOT
}

resource "aws_cloudfront_function" "rewrite_urls" {
  name    = "rewrite-urls-function"
  runtime = "cloudfront-js-2.0"
  code    = local.rewrite_function_code
  # Automatically publish the function to the LIVE stage
  publish = true

  # Learn more about the aws_cloudfront_function resource in the [Terraform Registry](https://registry.terraform.io)
}


# # Wait for the ACM certificate to be validated # TODO: ????
# resource "aws_acm_certificate_validation" "test" {
#   certificate_arn         = aws_acm_certificate.viewer_certificate.arn
#   # Use the FQDNs from the cloudflare_record resources to ensure Terraform waits for their creation
#   validation_record_fqdns = [for record in cloudflare_dns_record.acm_validation : record.name] # TODO: is name ok instead of hostname?
#
#   # TODO: Optional, Set a longer timeout if DNS propagation is slow...
#   timeouts {
#     create = "10m"
#   }
# }