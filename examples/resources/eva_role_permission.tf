
resource "eva_role_permission" "example_permission" {
  role_id = 1
  scoped_functionalities = jsonencode([
    {
      Functionality     = "FunctionalityName",
      Scope             = 1
      RequiresElevation = false
    }
  ])
}
