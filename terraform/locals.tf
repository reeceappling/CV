locals {
  s3_origin_id = "cvSiteBucketOrigin"
  source_dir = "./public/"
  full_domain = "${ var.subdomain }.${ var.domain }"
  site_bucket_name = local.full_domain
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