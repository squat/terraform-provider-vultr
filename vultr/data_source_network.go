package vultr

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	descriptionRegex, descriptionRegexOk := d.GetOk("description_regex")

	if !filtersOk && !descriptionRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "description_regex")
	}

	networks, err := client.GetNetworks()
	if err != nil {
		return fmt.Errorf("Error getting networks: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredNetworks []lib.Network
		for _, network := range networks {
			m := structToMap(network)
			if filter.F(m) {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		networks = filteredNetworks
	}

	if descriptionRegexOk {
		var filteredNetworks []lib.Network
		r := regexp.MustCompile(descriptionRegex.(string))
		for _, network := range networks {
			if r.MatchString(network.Description) {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		networks = filteredNetworks
	}

	if len(networks) < 1 {
		return errors.New("The query for networks returned no results. Please modify the search criteria and try again")
	}

	if len(networks) > 1 {
		return fmt.Errorf("The query for networks returned %d results. Please make the search criteria more specific and try again", len(networks))
	}

	d.SetId(networks[0].ID)
	d.Set("cidr_block", networkToCIDR(&networks[0]))
	d.Set("description", networks[0].Description)
	d.Set("region_id", networks[0].RegionID)
	return nil
}
