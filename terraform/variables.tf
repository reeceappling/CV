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
  description = "domain name to be hosted on (also the name of the s3 bucket). If using multi-nested subdomains, this should be the whole domain, less the first subdomain (ex: for cv.reece.appli.ng this value should be appli.ng)"
}

variable "aws_region" {
  type        = string
  default = "us-east-1"
  description = "aws region to use"
}

variable "subdomain" {
  type        = string
  description = "subdomains on the domain. (ex: for cv.reece.appli.ng, this value should be cv.reece)"
}