
resource "eva_organization_unit" "example" {
  name          = "example"
  phone_number  = "+316666666"
  email_address = "email@domain.com"
  backend_id    = "some-backend-id"
  parent_id     = 1
  currency_id   = "EUR"
}
