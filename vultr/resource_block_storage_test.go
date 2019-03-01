package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrBlockStorage_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrBlockStorageConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size", "20"),
				),
			},
			{
				Config: testAccVultrBlockStorageConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrBlockStorageExists("vultr_block_storage.foo"),
					resource.TestCheckResourceAttr("vultr_block_storage.foo", "size", "30"),
				),
			},
		},
	})
}

func testAccCheckVultrBlockStorageDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_block_storage" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Block Storage ID is not set")
		}

		storageID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetBlockStorage(storageID)

		if err == nil {
			return fmt.Errorf("Block Storage not deleted")
		}

		return nil
	}
	return nil
}

func testAccCheckVultrBlockStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Block Storage ID is not set")
		}

		storageID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		_, err := client.GetBlockStorage(storageID)

		if err != nil {
			return fmt.Errorf("Error getting Block Storage (%s): %v", storageID, err)
		}

		return nil
	}
}

func testAccVultrBlockStorageConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		data "vultr_region" "has_block_storage" {
			filter {
				name   = "block_storage"
				values = ["true"]
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

		resource "vultr_instance" "foo" {
			name              = "foo-%d"
			region_id         = "${data.vultr_region.has_block_storage.id}"
			plan_id           = "${data.vultr_plan.starter.id}"
			os_id             = "${data.vultr_os.container_linux.id}"
		}

		resource "vultr_block_storage" "foo" {
			name      = "foo-%d"
			region_id = "${data.vultr_region.has_block_storage.id}"
			size      = 20
			instance  = "${vultr_instance.foo.id}"
		}
	`, rInt, rInt)
}

func testAccVultrBlockStorageConfigUpdate(rInt int) string {
	return fmt.Sprintf(`  
		data "vultr_region" "has_block_storage" {
			filter {
				name   = "block_storage"
				values = ["true"]
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

		resource "vultr_instance" "foo" {
			name              = "foo-%d"
			region_id         = "${data.vultr_region.has_block_storage.id}"
			plan_id           = "${data.vultr_plan.starter.id}"
			os_id             = "${data.vultr_os.container_linux.id}"
		}

		resource "vultr_block_storage" "foo" {
			name      = "foo-%d"
			region_id = "${data.vultr_region.has_block_storage.id}"
			size      = 30
			instance  = "${vultr_instance.foo.id}"
		}
	`, rInt, rInt)
}
