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
				Config: testAccEvaRoleResourceConfig("my role", 1, "my_role"),
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
				Config: testAccEvaRoleResourceConfig("another role", 2, "another_role"),
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

func testAccEvaRoleResourceConfig(roleName string, userType int64, code string) string {
	return fmt.Sprintf(`
resource "eva_role" "test" {
	name          = "%s"
	user_type     = %d
	code          = "%s"
}
`, roleName, userType, code)
}
