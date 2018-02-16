package vultr

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallGroupCreate,
		Read:   resourceFirewallGroupRead,
		Update: resourceFirewallGroupUpdate,
		Delete: resourceFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"instance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceFirewallGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Creating new firewall group")
	id, err := client.CreateFirewallGroup(d.Get("description").(string))
	if err != nil {
		return fmt.Errorf("Error creating firewall group: %v", err)
	}

	d.SetId(id)

	return resourceFirewallGroupRead(d, meta)
}

func resourceFirewallGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	firewallGroup, err := client.GetFirewallGroup(d.Id())
	if err != nil {
		if err.Error() == fmt.Sprintf("Firewall group with ID %v not found", d.Id()) {
			log.Printf("[WARN] Removing firewall group (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting firewall group (%s): %v", d.Id(), err)
	}

	d.Set("description", firewallGroup.Description)
	d.Set("instance_count", firewallGroup.InstanceCount)
	d.Set("max_rule_count", firewallGroup.MaxRuleCount)
	d.Set("rule_count", firewallGroup.RuleCount)

	return nil
}

func resourceFirewallGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	if d.HasChange("description") {
		log.Printf("[INFO] Updating firewall group (%s) description", d.Id())
		if err := client.SetFirewallGroupDescription(d.Id(), d.Get("description").(string)); err != nil {
			return fmt.Errorf("Error changing firewall group (%s) description to %q: %v", d.Id(), d.Get("description").(string), err)
		}
	}

	return resourceFirewallGroupRead(d, meta)
}

func resourceFirewallGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying firewall group (%s)", d.Id())

	if err := client.DeleteFirewallGroup(d.Id()); err != nil {
		return fmt.Errorf("Error destroying firewall group (%s): %v", d.Id(), err)
	}

	return nil
}
