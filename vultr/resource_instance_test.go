package vultr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrInstance_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrInstanceConfigBasic(rInt, rSSH),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrInstanceExists("vultr_instance.bar"),
					resource.TestCheckResourceAttr("vultr_instance.bar", "tag", "bar"),
					resource.TestCheckResourceAttr("vultr_instance.bar", "ipv6", "false"),
				),
			},
			{
				Config: testAccVultrInstanceConfigUpdate(rInt, rSSH),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrInstanceExists("vultr_instance.bar"),
					resource.TestCheckResourceAttr("vultr_instance.bar", "tag", "baz"),
					resource.TestCheckResourceAttr("vultr_instance.bar", "ipv6", "true"),
				),
			},
		},
	})
}

func testAccCheckVultrInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance ID is not set")
		}
		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetServer(id)
		if err != nil {
			if strings.HasPrefix(err.Error(), "Invalid server") {
				return nil
			}
			return fmt.Errorf("Error getting instance (%s): %s", id, err)
		}

		return fmt.Errorf("Instance still exists: %s", id)
	}
	return nil
}

func testAccCheckVultrInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance ID is not set")
		}
		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetServer(id)
		if err != nil {
			return fmt.Errorf("Error getting instance (%s): %s", id, err)
		}

		return nil
	}
}

func testAccVultrInstanceConfigBasic(rInt int, rSSH string) string {
	return fmt.Sprintf(`
		data "vultr_region" "silicon_valley" {
			filter {
				name   = "name"
				values = ["Silicon Valley"]
			}
		}

		data "vultr_os" "container_linux" {
			filter {
				name   = "family"
				values = ["coreos"]
			}
		}

		data "vultr_plan" "starter" {
			filter {
				name   = "price_per_month"
				values = ["5.00"]
			}

			filter {
				name   = "ram"
				values = ["1024"]
			}
		}

		resource "vultr_ssh_key" "foo" {
			name       = "foo"
			public_key = "%s"
		}

		resource "vultr_instance" "bar" {
			name              = "bar-%d"
			region_id         = "${data.vultr_region.silicon_valley.id}"
			plan_id           = "${data.vultr_plan.starter.id}"
			os_id             = "${data.vultr_os.container_linux.id}"
			ssh_key_ids       = ["${vultr_ssh_key.foo.id}"]
			tag               = "bar"
			ipv6			  = false
		}
   `, rSSH, rInt)
}

func testAccVultrInstanceConfigUpdate(rInt int, rSSH string) string {
	return fmt.Sprintf(`
		data "vultr_region" "silicon_valley" {
			filter {
				name   = "name"
				values = ["Silicon Valley"]
			}
		}

		data "vultr_os" "container_linux" {
			filter {
				name   = "family"
				values = ["coreos"]
			}
		}

		data "vultr_plan" "2gb" {
			filter {
				name   = "ram"
				values = ["2048"]
			}
		}

		resource "vultr_ssh_key" "foo" {
			name       = "foo"
			public_key = "%s"
		}

		resource "vultr_instance" "bar" {
			name              = "bar-%d"
			region_id         = "${data.vultr_region.silicon_valley.id}"
			plan_id           = "${data.vultr_plan.2gb.id}"
			os_id             = "${data.vultr_os.container_linux.id}"
			ssh_key_ids       = ["${vultr_ssh_key.foo.id}"]
			tag               = "baz"
			ipv6			  = true
		}
   `, rSSH, rInt)
}
