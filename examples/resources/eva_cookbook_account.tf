resource "eva_cookbook_account" "example" {
  name           = "My Example Cookbook Account"
  object_account = "123" // Object account
  booking_flags  = 1 // WithTaxInformation
  type           = 1 // GeneralLedger
}
