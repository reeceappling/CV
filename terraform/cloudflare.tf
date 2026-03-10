resource "cloudflare_dns_record" "cv_subdomain_cname" {
  zone_id = var.cloudflare_zone_id
  name = "${ var.subdomain }.${ var.domain }"
  ttl     = 1
  type = "CNAME"
  comment = "Subdomain CNAME record for CV hosted on S3"
  content = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  proxied = true
  # settings = { # TODO: del or reenable???
  #   ipv4_only = true
  #   ipv6_only = true
  # }
  # tags = ["owner:dns-team"]
}