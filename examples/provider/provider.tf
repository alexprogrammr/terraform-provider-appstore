# Configuration for the App Store Connect provider.
provider "appstore" {
  key_id      = "2X9R4HXF34"
  issuer_id   = "57246542-96fe-1a63-e053-0824d011072a"
  private_key = file("key.p8")
}
