resource "eva_role" "example" {
  name      = "My Example Role"
  user_type = 1
  code      = "my_example_role"
  scoped_functionalities = [
    {
      functionality      = "FunctionalityName"
      scope              = 1
      requires_elevation = false
    }
  ]
}
