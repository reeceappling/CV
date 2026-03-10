# Ensure bucket exists
resource "aws_s3_bucket" "cv_site_bucket" {
  bucket = "${ var.subdomain}.${ var.domain }"
}

# Enable static website hosting
resource "aws_s3_bucket_website_configuration" "website_configuration" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  index_document {
    suffix = "index.html" # TODO: ensure ok
  }
  error_document {
    key    = "error.html"
  }
}

# Control public access settings
resource "aws_s3_bucket_public_access_block" "website_bucket_public_access_block" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# Set object ownership to allow ACLs (required for public-read ACL)
resource "aws_s3_bucket_ownership_controls" "website_bucket_ownership" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}
# Define bucket acl
resource "aws_s3_bucket_acl" "public_read" {
  depends_on = [
    aws_s3_bucket_ownership_controls.website_bucket_ownership,
    aws_s3_bucket_public_access_block.website_bucket_public_access_block,
  ]
  bucket = aws_s3_bucket.cv_site_bucket.id
  acl    = "public-read"
}

# Define a bucket policy to allow public read access
resource "aws_s3_bucket_policy" "public_bucket_policy" {
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

# TODO: bucket retention rules/backups (no backups)

# Upload all files in public directory to bucket
resource "aws_s3_object" "website_files" {
  for_each = fileset(local.source_dir, "**")
  bucket = aws_s3_bucket.cv_site_bucket.id
  key    = each.value
  source = "${local.source_dir}${each.value}"
  acl = "public-read"
  etag   = filemd5("${local.source_dir}/${each.value}") # Etag ensures updates are detected
  content_type = lookup(local.mime_types, regex("\\.([^.]+)$", each.value)[0], "application/octet-stream")
  //content_type = lookup(local.mime_types, regex("\\.([^.]+)$", each.value), "application/octet-stream") # TODO; NOT WORKING
}