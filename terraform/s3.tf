# Ensure bucket exists
resource "aws_s3_bucket" "cv_site_bucket" {
  bucket = local.site_bucket_name
  # TODO: tags???
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

# Set object ownership to allow ACLs (required for public-read ACL)
resource "aws_s3_bucket_ownership_controls" "website_bucket_ownership" { # TODO: might be unnecessary
  bucket = aws_s3_bucket.cv_site_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "b_acl" {
  depends_on = [aws_s3_bucket_ownership_controls.website_bucket_ownership] # TODO: ok?
  bucket = aws_s3_bucket.cv_site_bucket.id
  acl    = "private"
}

data "aws_iam_policy_document" "site_bucket" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.cv_site_bucket.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.cv_site_bucket.iam_arn]
    }
  }
}

resource "aws_s3_bucket_policy" "cv_site" {
  bucket = aws_s3_bucket.cv_site_bucket.id
  policy = data.aws_iam_policy_document.site_bucket.json
}

# TODO: bucket retention rules/backups (no backups)

# Upload all files in public directory to bucket
resource "aws_s3_object" "website_files" {
  for_each = fileset(local.source_dir, "**")
  bucket = aws_s3_bucket.cv_site_bucket.id
  key    = each.value
  source = "${local.source_dir}${each.value}"
  # acl = "public-read"
  etag   = filemd5("${local.source_dir}/${each.value}") # Etag ensures updates are detected
  //content_type = lookup(local.mime_types, regex("\\.([^.]+)$", each.value)[0], "application/octet-stream")
  content_type = lookup(local.mime_types, reverse(split(".", each.value))[0], "application/octet-stream")
  # TODO: lookup(local.mime_types, regex("\\.([^.]+)$", each.value)[0], "text/html")
}