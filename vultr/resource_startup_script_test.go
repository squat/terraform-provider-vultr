package vultr

import (
	"fmt"
	"testing"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrStartupScript_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrStartupScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrStartupScriptConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrStartupScriptExists("vultr_startup_script.foo"),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "type", "pxe"),
				),
			},
			{
				Config: testAccVultrStartupScriptConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrStartupScriptExists("vultr_startup_script.foo"),
					resource.TestCheckResourceAttr("vultr_startup_script.foo", "type", "boot"),
				),
			},
		},
	})
}

func testAccCheckVultrStartupScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_startup_script" {
			continue
		}

		scriptID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		script, err := client.GetStartupScript(scriptID)
		if err != nil {
			return fmt.Errorf("Error getting startup script (%s): %s", scriptID, err)
		}

		if script != (lib.StartupScript{}) {
			return fmt.Errorf("Startup script still exists: %s", scriptID)
		}

		return nil
	}
	return nil
}

func testAccCheckVultrStartupScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Script ID is not set")
		}

		scriptID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		script, err := client.GetStartupScript(scriptID)

		if err != nil {
			return fmt.Errorf("Error getting startup script (%s): %s", scriptID, err)
		}

		if script == (lib.StartupScript{}) {
			return fmt.Errorf("Startup script not found ID: %s", scriptID)
		}

		return nil
	}
}

func testAccVultrStartupScriptConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "foo-%d"
			type = "pxe"
			content = "#!/bin/bash\necho hello world > /root/hello"
		}
	`, rInt)
}

func testAccVultrStartupScriptConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_startup_script" "foo" {
			name = "foo-%d"
			type = "boot"
			content = "#!/bin/bash\necho hello world > /root/hello"
		}
	`, rInt)
}
