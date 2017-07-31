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

	_, cidrBlock, _ := net.ParseCIDR(d.Get("cidr_block").(string))
	firewallGroupID := d.Get("firewall_group_id").(string)
	fromPort := d.Get("from_port").(int)
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
	id, err := client.CreateFirewallRule(firewallGroupID, protocol, port, cidrBlock)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%d", firewallGroupID, id))

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idParts := strings.Split(d.Id(), "/")
	firewallGroupID := idParts[0]
	id, _ := strconv.Atoi(idParts[1])

	firewallRules, err := client.GetFirewallRules(firewallGroupID)
	if err != nil {
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
	d.Set("protocol", firewallRule.Protocol)
	ports := strings.Split(firewallRule.Port, ":")
	d.Set("from_port", ports[0])
	if len(ports) == 2 {
		d.Set("to_port", ports[1])
		return nil
	}
	d.Set("to_port", ports[0])

	return nil
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idParts := strings.Split(d.Id(), "/")
	firewallGroupID := idParts[0]
	id, _ := strconv.Atoi(idParts[1])

	log.Printf("[INFO] Destroying firewall rule (%s)", d.Id())

	if err := client.DeleteFirewallRule(id, firewallGroupID); err != nil {
		return fmt.Errorf("Error destroying firewall group (%s): %v", d.Id(), err)
	}

	return nil
}
