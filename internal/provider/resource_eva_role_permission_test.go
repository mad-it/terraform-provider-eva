package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaRolePermissionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccEvaRolePermissionsResourceConfig("AccessPrivacyDataRequests", 1, "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_role_permissions.test", "scoped_functionalities", "[{\"Functionality\":\"AccessPrivacyDataRequests\",\"RequiresElevation\":false,\"Scope\":1}]"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "eva_role_permissions.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccEvaRolePermissionsResourceConfig("InfrastructureProxy", 2, "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_role_permissions.test", "scoped_functionalities", "[{\"Functionality\":\"InfrastructureProxy\",\"RequiresElevation\":true,\"Scope\":2}]"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEvaRolePermissionsResourceConfig(functionality string, scope int64, requiresElevation string) string {
	return fmt.Sprintf(`
resource "eva_role" "test" {
	name          = "My role"
	user_type     = 1
	code          = "my_role"
}

resource "eva_role_permissions" "test" {
	role_id          		   = eva_role.test.id
	scoped_functionalities     = jsonencode([
		{
			Functionality     = "%s",
			Scope             = %d,
			RequiresElevation = %s
		}
	])
}
`, functionality, scope, requiresElevation)
}
