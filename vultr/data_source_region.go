package vultr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRegionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"block_storage": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"continent": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ddos_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	regions, err := client.GetRegions()
	if err != nil {
		return fmt.Errorf("Error getting regions: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredRegions []lib.Region
		for _, region := range regions {
			m := structToMap(region)
			if filter.F(m) {
				filteredRegions = append(filteredRegions, region)
			}
		}
		regions = filteredRegions
	}

	if nameRegexOk {
		var filteredRegions []lib.Region
		r := regexp.MustCompile(nameRegex.(string))
		for _, region := range regions {
			if r.MatchString(region.Name) {
				filteredRegions = append(filteredRegions, region)
			}
		}
		regions = filteredRegions
	}

	if len(regions) < 1 {
		return errors.New("The query for regions returned no results. Please modify the search criteria and try again")
	}

	if len(regions) > 1 {
		return fmt.Errorf("The query for regions returned %d results. Please make the search criteria more specific and try again", len(regions))
	}

	d.SetId(strconv.Itoa(regions[0].ID))
	d.Set("block_storage", regions[0].BlockStorage)
	d.Set("code", regions[0].Code)
	d.Set("continent", regions[0].Continent)
	d.Set("country", regions[0].Country)
	d.Set("ddos_protection", regions[0].Ddos)
	d.Set("name", regions[0].Name)
	d.Set("state", regions[0].State)
	return nil
}
