terraform {
  backend "s3" {
    // bucket gets overridden by pipeline
    bucket = var.domain
    key    = "tfState/cv-site.tfstate"
    region = "us-east-2"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.58.0"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
      version = "5.18.0"
    }
  }
}

# Configure Providers
provider "aws" {
  region = "us-east-2"
}
provider "cloudflare" {
  api_token = var.cloudflare_api_token
}


# Output the website endpoint URL
output "website_s3_endpoint" {
  value = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  description = "The S3 static website endpoint URL (HTTP only) (not through cloudflare)"
}
# Output the website endpoint URL for access through cloudflare
output "website_cloudflare_endpoint" {
  value = "${var.subdomain}.${var.domain}"
  description = "The https endpoint that goes through cloudflare"
}

# TODO: enable https
# TODO:

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