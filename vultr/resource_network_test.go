package vultr

import (
	"fmt"
	"testing"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrNetwork_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrNetworkConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrNetworkExists("vultr_network.foo"),
					resource.TestCheckResourceAttr("vultr_network.foo", "description", fmt.Sprintf("foo_%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckVultrNetworkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_network" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network ID is not set")
		}

		networkID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		networks, err := client.GetNetworks()
		if err != nil {
			return fmt.Errorf("Error getting network (%s): %v", networkID, err)
		}

		var network *lib.Network
		for _, n := range networks {
			if n.ID == networkID {
				network = &n
				break
			}
		}

		if network != nil {
			return fmt.Errorf("Network not deleted (%s)", networkID)
		}

		return nil
	}
	return nil
}

func testAccCheckVultrNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network ID is not set")
		}

		networkID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		networks, err := client.GetNetworks()
		if err != nil {
			return fmt.Errorf("Error getting network (%s): %v", networkID, err)
		}

		var network *lib.Network
		for _, n := range networks {
			if n.ID == networkID {
				network = &n
				break
			}
		}

		if network == nil {
			return fmt.Errorf("Network does not exist (%s)", networkID)
		}

		return nil
	}
}

func testAccVultrNetworkConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		data "vultr_region" "frankfurt" {
			filter {
				name   = "name"
				values = ["Frankfurt"]
			}
		}

		resource "vultr_network" "foo" {
			cidr_block  = "${cidrsubnet("192.168.0.0/23", 1, 0)}"
			description = "foo_%d"
			region_id   = "${data.vultr_region.frankfurt.id}"
		}
	`, rInt)
}
