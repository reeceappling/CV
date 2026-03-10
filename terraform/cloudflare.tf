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
  ttl     = 60 # Minimum value (seconds) for unproxied
  type    = each.value.type
  comment = "Certificate validation records for CV hosted on S3"
  content = trimsuffix(each.value.record, ".")
  proxied = false
  # tags = ["component:cv-site"]
}

resource "cloudflare_dns_record" "cv_subdomain_cname" {
  zone_id = var.cloudflare_zone_id
  name = local.domain_name
  type    = "CNAME"
  comment = "Subdomain CNAME record for CV hosted on S3"
  content = aws_cloudfront_distribution.s3_distribution.domain_name
  proxied = true
  ttl     = 1 # Must be 1 for proxied
  # tags = ["component:cv-site"]
}