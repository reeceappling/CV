variable "aws_region" {
  type        = string
  default = "us-east-1"
  description = "aws region to use"
}

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
  default = "appli.ng"
  description = "domain name to be hosted on. If using multi-nested subdomains, this should be only the TLD (ex: for cv.reece.appli.ng this value should be appli.ng)"
}

variable "subdomain" {
  type        = string
  default = "cv"
  description = "subdomains on the domain. (ex: for cv.reece.appli.ng, this value should be cv.reece)"
}
