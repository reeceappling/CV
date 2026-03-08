# Configure the Cloudflare Provider
terraform {
  required_providers {
    cloudflare = {
      source = "cloudflare/cloudflare"
      version = "5.18.0"
    }
  }
}

provider "aws" { region = "us-east-1" }

resource "aws_s3_bucket" "cv_site_bucket" {
  bucket = var.domain
}

# The provider block authenticates with Cloudflare
provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

# Create a CNAME record
resource "cloudflare_record" "www_cname" {
  zone_id = var.cloudflare_zone_id
  name    = "www"
  content = var.domain
  type    = "CNAME"
  proxied = true            # Set to true to enable Cloudflare proxy (orange cloud) # TODO; ?
  ttl     = 3600            # Optional: Time to live in seconds, "1" means automatic # TODO; ?
}

# Control public access settings
resource "aws_s3_bucket_public_access_block" "website_bucket_public_access_block" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# Enable static website hosting
resource "aws_s3_bucket_website_configuration" "website_configuration" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  index_document {
    suffix = "index.html" # TODO: ensure ok
  }
  error_document {
    suffix = ""
    key    = "error.html"
  }
}

# Set object ownership to allow ACLs (required for public-read ACL)
resource "aws_s3_bucket_ownership_controls" "website_bucket_ownership" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

# Define a bucket policy to allow public read access
resource "aws_s3_bucket_policy" "website_bucket_policy" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "PublicReadGetObject"
        Effect = "Allow"
        Principal = "*"
        Action = "s3:GetObject"
        Resource = "${aws_s3_bucket.cv_site_bucket.arn}/*"
      }
    ]
  })
}

locals {
  source_dir = "./public/"
  # Helper local for content types
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".png"  = "image/png"
    ".jpg"  = "image/jpeg"
    ".svg"  = "image/svg+xml"
  }
}

resource "aws_s3_object" "website_files" {
  for_each = fileset(local.source_dir, "**")
  bucket = aws_s3_bucket.cv_site_bucket.id
  key    = each.value
  source = "${local.source_dir}/${each.value}"
  acl = "public-read" # TODO: unsure???
  etag   = filemd5("${local.source_dir}/${each.value}") # Etag ensures updates are detected
  content_type = lookup(local.mime_types, regex("\\.([^.]+)$", each.value), "application/octet-stream")
}

# Output the website endpoint URL
output "website_endpoint" {
  value = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  description = "The S3 static website endpoint URL (HTTP only)"
}

# # Upload files from the local directory to the S3 bucket
# resource "aws_s3_object" "directory_upload" {
#   for_each = fileset("./public/", "**")
#   bucket = aws_s3_bucket.cv_site_bucket.id
#   # S3 path
#   key = each.value
#   # Local file source
#   source = "./public/${each.value}"
#   # include an etag to ensure updates iff file content changes
#   etag = filemd5("./public/${each.value}")
# }