# Create a CNAME record
resource "cloudflare_record" "www_cname" {
  zone_id = var.cloudflare_zone_id
  name    = var.subdomain
  content = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  type    = "CNAME"
  proxied = true
  ttl     = 600            # Optional: Time to live in seconds, "1" means automatic # TODO; ?
}