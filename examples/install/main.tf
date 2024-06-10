terraform {
  required_providers {
    appstore = {
      source = "alexprogrammr/appstore"
    }
  }
}

provider "appstore" {
  key_id      = ""
  issuer_id   = ""
  private_key = file("")
}

data "appstore_apps" "apps" {}

output "out" {
  value = data.appstore_apps.apps
}
