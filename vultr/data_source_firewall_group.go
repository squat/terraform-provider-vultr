package vultr

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallGroupRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceFirewallGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	descriptionRegex, descriptionRegexOk := d.GetOk("description_regex")

	if !filtersOk && !descriptionRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "description_regex")
	}

	firewallGroups, err := client.GetFirewallGroups()
	if err != nil {
		return fmt.Errorf("Error getting firewall groups: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredFirewallGroups []lib.FirewallGroup
		for _, firewallGroup := range firewallGroups {
			m := structToMap(firewallGroup)
			if filter.F(m) {
				filteredFirewallGroups = append(filteredFirewallGroups, firewallGroup)
			}
		}
		firewallGroups = filteredFirewallGroups
	}

	if descriptionRegexOk {
		var filteredFirewallGroups []lib.FirewallGroup
		r := regexp.MustCompile(descriptionRegex.(string))
		for _, firewallGroup := range firewallGroups {
			if r.MatchString(firewallGroup.Description) {
				filteredFirewallGroups = append(filteredFirewallGroups, firewallGroup)
			}
		}
		firewallGroups = filteredFirewallGroups
	}

	if len(firewallGroups) < 1 {
		return errors.New("The query for firewall groups returned no results. Please modify the search criteria and try again")
	}

	if len(firewallGroups) > 1 {
		return fmt.Errorf("The query for firewall groups returned %d results. Please make the search criteria more specific and try again", len(firewallGroups))
	}

	d.SetId(firewallGroups[0].ID)
	d.Set("description", firewallGroups[0].Description)
	d.Set("instance_count", firewallGroups[0].InstanceCount)
	d.Set("max_rule_count", firewallGroups[0].MaxRuleCount)
	d.Set("rule_count", firewallGroups[0].RuleCount)
	return nil
}
