package vultr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePlanRead,

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
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Computed: true,
			},

			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	plans, err := client.GetPlans()
	if err != nil {
		return fmt.Errorf("Error getting plans: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredPlans []lib.Plan
		for _, plan := range plans {
			m := structToMap(plan)
			if filter.F(m) {
				filteredPlans = append(filteredPlans, plan)
			}
		}
		plans = filteredPlans
	}

	if nameRegexOk {
		var filteredPlans []lib.Plan
		r := regexp.MustCompile(nameRegex.(string))
		for _, plan := range plans {
			if r.MatchString(plan.Name) {
				filteredPlans = append(filteredPlans, plan)
			}
		}
		plans = filteredPlans
	}

	if len(plans) < 1 {
		return errors.New("The query for plans returned no results. Please modify the search criteria and try again")
	}

	if len(plans) > 1 {
		return fmt.Errorf("The query for plans returned %d results. Please make the search criteria more specific and try again", len(plans))
	}

	d.SetId(strconv.Itoa(plans[0].ID))
	d.Set("available_locations", plans[0].Regions)
	d.Set("bandwidth", plans[0].Bandwidth)
	d.Set("disk", plans[0].Disk)
	d.Set("name", plans[0].Name)
	d.Set("price_per_month", plans[0].Price)
	d.Set("ram", plans[0].RAM)
	d.Set("vcpu_count", plans[0].VCpus)
	return nil
}
