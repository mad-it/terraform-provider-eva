
resource "eva_role_permission" "example_permission" {
  role_id = 1
  scoped_functionalities = [
    {
      functionality      = "FunctionalityName"
      scope              = 1
      requires_elevation = false
    }
  ]
}
