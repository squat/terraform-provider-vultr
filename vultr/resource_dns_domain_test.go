package vultr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrDNSDomain_basic(t *testing.T) {
	t.Parallel()

	rStr := acctest.RandomWithPrefix("tf-test-")
	domain := rStr + ".dev"
	ip := "10.0.0.1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrDNSDomainConfigBasic(domain, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrDNSDomainExists("vultr_dns_domain.foo"),
					resource.TestCheckResourceAttr("vultr_dns_domain.foo", "domain", domain),
					resource.TestCheckResourceAttr("vultr_dns_domain.foo", "ip", ip),
				),
			},
		},
	})
}

func testAccCheckVultrDNSDomainDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_dns_domain" {
			continue
		}

		client := testAccProvider.Meta().(*Client)

		dnsDomains, err := client.GetDNSDomains()
		if err != nil {
			return fmt.Errorf("Error getting DNS domains: %v", err)
		}

		for i := range dnsDomains {
			if dnsDomains[i].Domain == rs.Primary.ID {
				return fmt.Errorf("DNS Domain still exists: %v", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckVultrDNSDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client := testAccProvider.Meta().(*Client)

		dnsDomains, err := client.GetDNSDomains()
		if err != nil {
			return fmt.Errorf("Error getting DNS domains: %v", err)
		}

		for i := range dnsDomains {
			if dnsDomains[i].Domain == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("DNS Domain does not exist: %v", rs.Primary.ID)
	}
}

func testAccVultrDNSDomainConfigBasic(domain string, ip string) string {
	return fmt.Sprintf(`
		resource "vultr_dns_domain" "foo" {
			domain = "%s"
			ip     = "%s"
		}
   `, domain, ip)
}
