package vultr

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceStartupScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStartupScriptRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStartupScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	startupScripts, err := client.GetStartupScripts()
	if err != nil {
		return fmt.Errorf("Error getting startup scripts: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredStartupScripts []lib.StartupScript
		for _, startupScript := range startupScripts {
			m := structToMap(startupScript)
			if filter.F(m) {
				filteredStartupScripts = append(filteredStartupScripts, startupScript)
			}
		}
		startupScripts = filteredStartupScripts
	}

	if nameRegexOk {
		var filteredStartupScripts []lib.StartupScript
		r := regexp.MustCompile(nameRegex.(string))
		for _, startupScript := range startupScripts {
			if r.MatchString(startupScript.Name) {
				filteredStartupScripts = append(filteredStartupScripts, startupScript)
			}
		}
		startupScripts = filteredStartupScripts
	}

	if len(startupScripts) < 1 {
		return errors.New("The query for startup scripts returned no results. Please modify the search criteria and try again")
	}

	if len(startupScripts) > 1 {
		return fmt.Errorf("The query for startup scripts returned %d results. Please make the search criteria more specific and try again", len(startupScripts))
	}

	d.SetId(startupScripts[0].ID)
	d.Set("content", startupScripts[0].Content)
	d.Set("name", startupScripts[0].Name)
	d.Set("type", startupScripts[0].Type)
	return nil
}
