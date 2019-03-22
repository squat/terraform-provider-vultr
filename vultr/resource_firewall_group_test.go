package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrFirewallGroup_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallGroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.foo"),
					resource.TestCheckResourceAttr("vultr_firewall_group.foo", "description", fmt.Sprintf("foo-%d", rInt)),
				),
			},
			{
				Config: testAccVultrFirewallGroupConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallGroupExists("vultr_firewall_group.foo"),
					resource.TestCheckResourceAttr("vultr_firewall_group.foo", "description", fmt.Sprintf("bar-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckVultrFirewallGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_firewall_group" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Firewall Group ID is not set")
		}

		fwgID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetFirewallGroup(fwgID)

		if err == nil {
			return fmt.Errorf("Firewall Group not deleted")
		}

		return nil
	}
	return nil
}

func testAccCheckVultrFirewallGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Firewall Group ID is not set")
		}

		fwgID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetFirewallGroup(fwgID)

		if err != nil {
			return fmt.Errorf("Error getting Firewall Group (%s): %v", fwgID, err)
		}

		return nil
	}
}

func testAccVultrFirewallGroupConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "foo" {
			description = "foo-%d"
		}
	`, rInt)
}

func testAccVultrFirewallGroupConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "foo" {
			description = "bar-%d"
		}
	`, rInt)
}
