package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEvaRoleResourceConfig(
					roleConfig{
						name:     "my role",
						userType: 1,
						code:     "my_role",
					},
					roleScopedFunctionalityConfig{
						functionality:      "Some functionality",
						scope:              1,
						requires_elevation: "false",
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_role.test", "name", "my role"),
					resource.TestCheckResourceAttr("eva_role.test", "user_type", "1"),
					resource.TestCheckResourceAttr("eva_role.test", "code", "my_role"),
				),
			},
			// // ImportState testing
			// {
			// 	ResourceName:      "eva_role.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: testAccEvaRoleResourceConfig(
					roleConfig{
						name:     "another role",
						userType: 2,
						code:     "another_role",
					},
					roleScopedFunctionalityConfig{
						functionality:      "Another functionality",
						scope:              2,
						requires_elevation: "true",
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_role.test", "name", "another role"),
					resource.TestCheckResourceAttr("eva_role.test", "user_type", "2"),
					resource.TestCheckResourceAttr("eva_role.test", "code", "another_role"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

type roleConfig struct {
	name     string
	userType int64
	code     string
}

type roleScopedFunctionalityConfig struct {
	functionality      string
	scope              int64
	requires_elevation string
}

func testAccEvaRoleResourceConfig(roleConfig roleConfig, permissionConfig roleScopedFunctionalityConfig) string {
	return fmt.Sprintf(`
resource "eva_role" "test" {
	name                   = "%s"
	user_type              = %d
	code                   = "%s"
	scoped_functionalities = [
		{
			functionality      = "%s"
			scope              = %d
			requires_elevation = %s
		}
	]
}`, roleConfig.name, roleConfig.userType, roleConfig.code, permissionConfig.functionality, permissionConfig.scope, permissionConfig.requires_elevation)
}
