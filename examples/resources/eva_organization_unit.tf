
resource "eva_organization_unit" "example" {
  name          = "example"
  phone_number  = "+316666666"
  email_address = "email@domain.com"
  backend_id    = "some-backend-id"
  parent_id     = 1
  currency_id   = "EUR"
  type          = 8
  address = {
    address1     = "address 1"
    address2     = "address 2"
    house_number = "22"
    zip_code     = "30254 AD"
    city         = "Amsterdam"
    country_id   = "NL"
    latitude     = 41.27578565
    longitude    = -8.28065198
  }
}
