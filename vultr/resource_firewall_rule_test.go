package vultr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVultrFirewallRule_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVultrFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrFirewallRuleConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrFirewallRuleExists("vultr_firewall_rule.bar"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.bar", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vultr_firewall_rule.bar", "cidr_block", "0.0.0.0/0"),
				),
			},
		},
	})
}

func testAccCheckVultrFirewallRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_firewall_rule" {
			continue
		}

		firewallGroupID, id, err := parseStringSlashInt(rs.Primary.ID, "firewall rule ID", "firewall-group-id", "firewall-rule-number")
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*Client)
		firewallRules, err := client.GetFirewallRules(firewallGroupID)
		if err != nil {
			if strings.HasPrefix(err.Error(), "Invalid firewall group") {
				// firewall group and rule were destroyed by Terraform cleanup
				return nil
			}
			return fmt.Errorf("Error getting firewall rule (%s): %v", firewallGroupID, err)
		}

		for _, f := range firewallRules {
			if f.RuleNumber == id {
				return fmt.Errorf("Firewall rule was not deleted")
			}
		}

		return nil
	}
	return nil
}

func testAccCheckVultrFirewallRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		firewallGroupID, id, err := parseStringSlashInt(rs.Primary.ID, "firewall rule ID", "firewall-group-id", "firewall-rule-number")
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*Client)
		firewallRules, err := client.GetFirewallRules(firewallGroupID)
		if err != nil {
			return fmt.Errorf("Error retrieving firewall rules: %v", err)
		}

		for _, f := range firewallRules {
			if f.RuleNumber == id {
				return nil
			}
		}

		return fmt.Errorf("Firewall rule does not exist")
	}
}

func testAccVultrFirewallRuleConfigBasic(rInt int) string {
	return fmt.Sprintf(`
		resource "vultr_firewall_group" "foo" {
			description = "foo-%d"
		}

		resource "vultr_firewall_rule" "bar" {
			firewall_group_id = "${vultr_firewall_group.foo.id}"
			cidr_block        = "0.0.0.0/0"
			protocol          = "tcp"
			from_port         = 22
			to_port           = 22
		}
	`, rInt)
}

func TestSplitFirewallRule(t *testing.T) {
	cases := []struct {
		portRange string
		from      int
		to        int
		err       bool
	}{
		{
			portRange: "",
			from:      0,
			to:        0,
			err:       false,
		},
		{
			portRange: ":",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "-",
			from:      0,
			to:        0,
			err:       false,
		},
		{
			portRange: " - ",
			from:      0,
			to:        0,
			err:       false,
		},
		{
			portRange: "22",
			from:      22,
			to:        22,
			err:       false,
		},
		{
			portRange: "foo",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "22:23",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "22 - 23",
			from:      22,
			to:        23,
			err:       false,
		},
		{
			portRange: "80-81",
			from:      80,
			to:        81,
			err:       false,
		},
	}

	for i, c := range cases {
		from, to, err := splitFirewallRule(c.portRange)
		if (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
		if from != c.from || to != c.to {
			t.Errorf("test case %d: expected range %d:%d, got %d:%d", i, c.from, c.to, from, to)
		}
	}
}
