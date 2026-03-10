# Create a CNAME record
# resource "cloudflare_record" "cv_subdomain_cname" {
#   zone_id = var.cloudflare_zone_id
#   name    = var.subdomain
#   content = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
#   type    = "CNAME"
#   proxied = true
#   ttl     = 600            # Optional: Time to live in seconds, "1" means automatic # TODO; was 3600 (30min), set to 10min for testing purposes
# }

resource "cloudflare_dns_record" "cv_subdomain_cname" {
  zone_id = var.cloudflare_zone_id
  name = "${ var.subdomain }.${ var.domain }"
  ttl     = 600            # Optional: Time to live in seconds, "1" means automatic # TODO; was 3600 (30min), set to 10min for testing purposes
  type = "CNAME"
  comment = "Subdomain CNAME record for CV hosted on S3"
  content = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  proxied = true
  # settings = {
  #   ipv4_only = true
  #   ipv6_only = true
  # }
  # tags = ["owner:dns-team"]
}