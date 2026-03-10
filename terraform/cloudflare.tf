resource "cloudflare_dns_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.viewer_certificate.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  zone_id = var.cloudflare_zone_id
  name = each.value.name
  ttl     = 60
  type    = each.value.type
  comment = "Subdomain CNAME record for CV hosted on S3" # TODO: fix
  # Use trimsuffix to remove the trailing dot that AWS sometimes adds, which Cloudflare doesn't require
  value   = trimsuffix(each.value.record, ".")
  # content = trimsuffix(each.value.record, ".")
  proxied = false
  # settings = { # TODO: del or reenable???
  #   ipv4_only = true
  #   ipv6_only = true
  # }
  # tags = ["owner:dns-team"]
}

resource "cloudflare_dns_record" "cv_subdomain_cname" {
  zone_id = var.cloudflare_zone_id
  name = local.domain_name
  type    = "CNAME"
  comment = "Subdomain CNAME record for CV hosted on S3"
  # TODO: value = aws_cloudfront_distribution.cloudfront_distribution.domain_name
  content = aws_cloudfront_distribution.s3_distribution.domain_name
  proxied = true
  ttl     = 1 # Must be 1 for proxied
  # settings = { # TODO: del or reenable???
  #   ipv4_only = true
  #   ipv6_only = true
  # }
  # tags = ["owner:dns-team"]
}