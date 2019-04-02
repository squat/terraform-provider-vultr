package vultr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrReservedIP_basic(t *testing.T) {
	t.Parallel()

	rName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfigBasicIPv4(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "name", rName),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "type", "v4"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "cidr"),
				),
			},
			{
				Config: testAccVultrReservedIPConfigBasicIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "name", rName),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "type", "v6"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "cidr"),
				),
			},
		},
	})
}

func testAccCheckVultrReservedIPDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_reserved_ip" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Reserved IP ID is not set")
		}

		client := testAccProvider.Meta().(*Client)

		_, err := client.GetReservedIP(rs.Primary.ID)
		if err != nil {
			if strings.HasPrefix(err.Error(), fmt.Sprintf("IP with ID %v not found", rs.Primary.ID)) {
				return nil
			}
			return fmt.Errorf("Error getting reserved ip (%s): %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("Reserved IP still exists: (%s)", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVultrReservedIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Reserved IP ID is not set")
		}

		client := testAccProvider.Meta().(*Client)

		_, err := client.GetReservedIP(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting reserved ip (%s): %v", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccVultrReservedIPConfigBasicIPv4(name string) string {
	return fmt.Sprintf(`
	data "vultr_region" "silicon_valley" {
		filter {
		  name   = "name"
		  values = ["Silicon Valley"]
		}
	}
	resource "vultr_reserved_ip" "foo" {
		name        = "%s"
		region_id   = "${data.vultr_region.silicon_valley.id}"
		type		= "v4"
	}
   `, name)
}

func testAccVultrReservedIPConfigBasicIPv6(name string) string {
	return fmt.Sprintf(`
	data "vultr_region" "silicon_valley" {
		filter {
		  name   = "name"
		  values = ["Silicon Valley"]
		}
	}
	resource "vultr_reserved_ip" "foo" {
		name        = "%s"
		region_id   = "${data.vultr_region.silicon_valley.id}"
		type		= "v6"
	}
   `, name)
}
