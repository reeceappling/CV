terraform {
  backend "s3" {
    // bucket and region get overridden by pipeline. Ex: terraform init -backend-config='bucket=stateBucketName'
    bucket = "overwritten-by-pipeline"
    key    = "tfState/cv-site.tfstate"
    region = "us-east-1"
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
  region = var.aws_region
}
provider "cloudflare" {
  api_token = var.cloudflare_api_token
}
# Provider for us-east-1 region (required for CloudFront certs)
provider "aws" {
  region = "us-east-1"
  alias  = "us_east_1"
}


# Output the website endpoint URL
output "website_s3_endpoint" {
  value = aws_s3_bucket_website_configuration.website_configuration.website_endpoint
  description = "The S3 static website endpoint URL (HTTP only) (not through cloudflare)"
}
# Output the website endpoint URL for access through cloudflare
output "website_cloudflare_endpoint_and_bucket_name" {
  value = "${var.subdomain}.${var.domain}"
  description = "The https endpoint that goes through cloudflare"
}

# TODO: enable https?
# TODO: