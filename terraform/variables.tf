variable "cloudflare_api_token" {
  type        = string
  sensitive   = true
  description = "Cloudflare API token."
}

variable "cloudflare_zone_id" {
  type        = string
  description = "Zone ID of your Cloudflare domain."
}

variable "domain" {
  type        = string
  description = "domain name to be hosted on (also the name of the s3 bucket)"
}