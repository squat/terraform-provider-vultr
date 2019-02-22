package vultr

import (
	"fmt"
	"testing"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrSSHKey_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rSSH, _, err := acctest.RandSSHKeyPair("foobar")

	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrSSHKeyConfigBasic(rInt, rSSH),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrSSHKeyExists("vultr_ssh_key.foo"),
					resource.TestCheckResourceAttr("vultr_ssh_key.foo", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
			{
				Config: testAccVultrSSHKeyConfigUpdate(rInt, rSSH),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrSSHKeyExists("vultr_ssh_key.foo"),
					resource.TestCheckResourceAttr("vultr_ssh_key.foo", "name", fmt.Sprintf("bar-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckVultrSSHKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_ssh_key" {
			continue
		}

		keyID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		keys, err := client.GetSSHKeys()
		if err != nil {
			return fmt.Errorf("Error getting SSH keys: %s", err)
		}

		var key *lib.SSHKey
		for i := range keys {
			if keys[i].ID == keyID {
				key = &keys[i]
				break
			}
		}

		if key != nil {
			return fmt.Errorf("SSH Key still exists: %s", keyID)
		}
	}
	return nil
}

func testAccCheckVultrSSHKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SSH Key ID is not set")
		}

		keyID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)

		keys, err := client.GetSSHKeys()
		if err != nil {
			return fmt.Errorf("Error getting SSH keys: %s", err)
		}

		var key *lib.SSHKey
		for i := range keys {
			if keys[i].ID == keyID {
				key = &keys[i]
				break
			}
		}

		if key == nil {
			return fmt.Errorf("SSH Key does not exist: %s", keyID)
		}

		return nil
	}
}

func testAccVultrSSHKeyConfigBasic(rInt int, rSSH string) string {
	return fmt.Sprintf(`
		resource "vultr_ssh_key" "foo" {
			name       = "foo-%d"
			public_key = "%s"
		}
	`, rInt, rSSH)
}

func testAccVultrSSHKeyConfigUpdate(rInt int, rSSH string) string {
	return fmt.Sprintf(`
		resource "vultr_ssh_key" "foo" {
			name       = "bar-%d"
			public_key = "%s"
		}
	`, rInt, rSSH)
}
