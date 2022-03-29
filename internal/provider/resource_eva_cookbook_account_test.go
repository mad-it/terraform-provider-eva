package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaCookbookAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEvaCookbookAccountResourceConfig("cookbook account", "123", "1", "1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "name", "cookbook account"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "object_account", "123"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "booking_flags", "1"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "type", "1"),
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
				Config: testAccEvaCookbookAccountResourceConfig("my account", "321", "2", "3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "name", "my account"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "object_account", "321"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "booking_flags", "2"),
					resource.TestCheckResourceAttr("eva_cookbook_account.test", "type", "3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEvaCookbookAccountResourceConfig(name string, objectAccount string, bookingFlags string, accountType string) string {
	return fmt.Sprintf(`
resource "eva_cookbook_account" "test" {
	name                   = "%s"
	object_account         = "%s"
	booking_flags          = %s
	type                   = %s
}`, name, objectAccount, bookingFlags, accountType)
}
