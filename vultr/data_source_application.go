package vultr

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"deploy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"short_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"surcharge": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	applications, err := client.GetApplications()
	if err != nil {
		return fmt.Errorf("Error getting applications: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredApplications []lib.Application
		for _, application := range applications {
			m := structToMap(application)
			if filter.F(m) {
				filteredApplications = append(filteredApplications, application)
			}
		}
		applications = filteredApplications
	}

	if nameRegexOk {
		var filteredApplications []lib.Application
		r := regexp.MustCompile(nameRegex.(string))
		for _, application := range applications {
			if r.MatchString(application.Name) {
				filteredApplications = append(filteredApplications, application)
			}
		}
		applications = filteredApplications
	}

	if len(applications) < 1 {
		return errors.New("The query for applications returned no results. Please modify the search criteria and try again")
	}

	if len(applications) > 1 {
		return fmt.Errorf("The query for applications returned %d results. Please make the search criteria more specific and try again", len(applications))
	}

	d.SetId(applications[0].ID)
	d.Set("deploy_name", applications[0].DeployName)
	d.Set("name", applications[0].Name)
	d.Set("short_name", applications[0].ShortName)
	d.Set("surcharge", applications[0].Surcharge)
	return nil
}
