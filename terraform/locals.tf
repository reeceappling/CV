locals {
  s3_origin_id = "cvSiteBucketOrigin"
  source_dir = "./public/"
  domain_name = "${ var.subdomain }.${ var.domain }"
  site_bucket_name = local.domain_name
  # Helper local for content types
  mime_types = {
    "html" = "text/html"
    "css"  = "text/css"
    "js"   = "application/javascript"
    "png"  = "image/png"
    "jpg"  = "image/jpeg"
    "svg"  = "image/svg+xml"
    "webp"  = "image/webp"
  }
}