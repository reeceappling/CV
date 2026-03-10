locals {
  s3_origin_id = "cvSiteBucketOrigin"
  source_dir = "./public/"
  site_bucket_name = "${ var.subdomain}.${ var.domain }" # TODO: del?
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