
resource "eva_open_id_provider" "provider" {
  name                = "provider"
  enabled             = true
  base_url            = "https://some-oauth-url.com"
  client_id           = "client-id"
  first_name_claim    = "given_name"
  last_name_claim     = "family_name"
  email_address_claim = "email"
  nickname_claim      = "name"
  user_type           = 1
  create_users        = true
  primary             = true
}
