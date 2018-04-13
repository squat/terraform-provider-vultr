package vultr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceBareMetalPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBareMetalPlanRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"available_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"price_per_month": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceBareMetalPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	bareMetalPlans, err := client.GetBareMetalPlans()
	if err != nil {
		return fmt.Errorf("Error getting bare metal plans: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredBareMetalPlans []lib.BareMetalPlan
		for _, bareMetalPlan := range bareMetalPlans {
			m := structToMap(bareMetalPlan)
			if filter.F(m) {
				filteredBareMetalPlans = append(filteredBareMetalPlans, bareMetalPlan)
			}
		}
		bareMetalPlans = filteredBareMetalPlans
	}

	if nameRegexOk {
		var filteredBareMetalPlans []lib.BareMetalPlan
		r := regexp.MustCompile(nameRegex.(string))
		for _, bareMetalPlan := range bareMetalPlans {
			if r.MatchString(bareMetalPlan.Name) {
				filteredBareMetalPlans = append(filteredBareMetalPlans, bareMetalPlan)
			}
		}
		bareMetalPlans = filteredBareMetalPlans
	}

	if len(bareMetalPlans) < 1 {
		return errors.New("The query for bare metal plans returned no results. Please modify the search criteria and try again")
	}

	if len(bareMetalPlans) > 1 {
		return fmt.Errorf("The query for bare metal plans returned %d results. Please make the search criteria more specific and try again", len(bareMetalPlans))
	}

	d.SetId(strconv.Itoa(bareMetalPlans[0].ID))
	d.Set("available_locations", bareMetalPlans[0].Regions)
	d.Set("bandwidth", bareMetalPlans[0].Bandwidth)
	d.Set("disk", bareMetalPlans[0].Disk)
	d.Set("name", bareMetalPlans[0].Name)
	d.Set("price_per_month", bareMetalPlans[0].Price)
	d.Set("ram", bareMetalPlans[0].RAM)
	d.Set("cpu_count", bareMetalPlans[0].CPUs)
	return nil
}
