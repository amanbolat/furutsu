terraform {
  required_version = ">= 0.10"

  backend "s3" {
    region = "us-west-2"
    bucket = "terraform.state.furutsu.amanbolat.com"
    key = "terraform"
    profile = "prod"
  }
}

variable "profile" {}

variable "region" {}

variable "domain" {
  default = "furutsu.amanbolat.com"
}

variable "deployer_acc" {}

variable "cert_arn" {}

data "aws_route53_zone" "domain_zone" {
  name = "amanbolat.com."
}

provider "aws" {
  region = var.region
  profile = var.profile
}

module "site-main" {
  source = "github.com/ringods/terraform-website-s3-cloudfront-route53//site-main"

  region = var.region
  domain = var.domain
  bucket_name = "${var.domain}-static"
  duplicate-content-penalty-secret = "dvyRK=3faA!31do30g0A#fa22ggjoaan4L%bL7Xvj"
  deployer = var.deployer_acc
  acm-certificate-arn = var.cert_arn
  not-found-response-path = "/index.html"
}

module "dns-alias" {
  source = "github.com/ringods/terraform-website-s3-cloudfront-route53//r53-alias"

  domain = var.domain
  target = module.site-main.website_cdn_hostname
  cdn_hosted_zone_id = module.site-main.website_cdn_zone_id
  route53_zone_id = data.aws_route53_zone.domain_zone.zone_id
}