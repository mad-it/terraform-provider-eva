package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaOrderLedgerTypeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEvaOrderLedgerTypeResourceConfig("pending", "pending order"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_order_ledger_type.test", "name", "pending"),
					resource.TestCheckResourceAttr("eva_order_ledger_type.test", "description", "pending order"),
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
				Config: testAccEvaCustomOrderStatusResourceConfig("completed", "completed order"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("eva_order_ledger_type.test", "name", "completed"),
					resource.TestCheckResourceAttr("eva_order_ledger_type.test", "description", "completed order"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEvaOrderLedgerTypeResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "eva_order_ledger_type" "test" {
	name                   = "%s"
	description            = "%s"
}`, name, description)
}
