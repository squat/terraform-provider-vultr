package vultr

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Delete: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"direction": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"firewall_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"from_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateFirewallRuleProtocol,
			},

			"to_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	_, cidrBlock, err := net.ParseCIDR(d.Get("cidr_block").(string))
	if err != nil {
		return fmt.Errorf("Error parsing %q for firewall rule: %v", "cidr_block", err)
	}
	firewallGroupID := d.Get("firewall_group_id").(string)
	fromPort := d.Get("from_port").(int)
	notes := d.Get("notes").(string)
	protocol := d.Get("protocol").(string)
	toPort := d.Get("to_port").(int)

	_, fok := d.GetOk("from_port")
	_, tok := d.GetOk("to_port")
	if fok != tok {
		return fmt.Errorf("Expected %q and %q to both be provided or both be empty", "from_port", "to_port")
	}

	if (protocol == "tcp" || protocol == "udp") && !fok {
		return fmt.Errorf("%q and %q are required for protocol of type %q", "from_port", "to_port", protocol)
	}

	var port string
	if fok {
		if fromPort != toPort {
			port = fmt.Sprintf("%d:%d", fromPort, toPort)
		} else {
			port = strconv.Itoa(fromPort)
		}
	}

	log.Printf("[INFO] Creating new firewall rule")
	id, err := client.CreateFirewallRule(firewallGroupID, protocol, port, cidrBlock, notes)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%d", firewallGroupID, id))

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	firewallGroupID, id, err := parseStringSlashInt(d.Id(), "firewall rule ID", "firewall-group-id", "firewall-rule-number")
	if err != nil {
		return err
	}

	firewallRules, err := client.GetFirewallRules(firewallGroupID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Invalid firewall group") {
			log.Printf("[WARN] Removing firewall rule (%s) because the group is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting firewall rule (%s): %v", d.Id(), err)
	}

	var firewallRule *lib.FirewallRule
	for _, f := range firewallRules {
		if f.RuleNumber == id {
			firewallRule = &f
			break
		}
	}

	if firewallRule == nil {
		log.Printf("[WARN] Removing firewall rule (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("action", firewallRule.Action)
	d.Set("cidr_block", firewallRule.Network.String())
	d.Set("direction", "in")
	d.Set("firewall_group_id", firewallGroupID)
	d.Set("notes", firewallRule.Notes)
	d.Set("protocol", firewallRule.Protocol)
	from, to, err := splitFirewallRule(firewallRule.Port)
	if err != nil {
		return fmt.Errorf("Error parsing port range for firewall rule (%s): %v", d.Id(), err)
	}
	d.Set("from_port", from)
	d.Set("to_port", to)

	return nil
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	firewallGroupID, id, err := parseStringSlashInt(d.Id(), "firewall rule ID", "firewall-group-id", "firewall-rule-number")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Destroying firewall rule (%s)", d.Id())

	if err := client.DeleteFirewallRule(id, firewallGroupID); err != nil {
		return fmt.Errorf("Error destroying firewall rule (%s): %v", d.Id(), err)
	}

	return nil
}

func splitFirewallRule(portRange string) (int, int, error) {
	if len(portRange) == 0 || strings.TrimSpace(portRange) == "-" {
		return 0, 0, nil
	}
	ports := strings.Split(portRange, "-")
	from, err := strconv.Atoi(strings.TrimSpace(ports[0]))
	if err != nil {
		return 0, 0, err
	}
	if len(ports) == 1 {
		return from, from, nil
	}
	to, err := strconv.Atoi(strings.TrimSpace(ports[1]))
	if err != nil {
		return 0, 0, err
	}
	return from, to, nil
}
