package vultr

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	keys, err := client.GetSSHKeys()
	if err != nil {
		return fmt.Errorf("Error getting SSH keys: %v", err)
	}

	filter := filterFromSet(filters.(*schema.Set))
	if filtersOk {
		var filteredKeys []lib.SSHKey
		for _, key := range keys {
			m := structToMap(key)
			if filter.F(m) {
				filteredKeys = append(filteredKeys, key)
			}
		}
		keys = filteredKeys
	}

	if nameRegexOk {
		var filteredKeys []lib.SSHKey
		r := regexp.MustCompile(nameRegex.(string))
		for _, key := range keys {
			if r.MatchString(key.Name) {
				filteredKeys = append(filteredKeys, key)
			}
		}
		keys = filteredKeys
	}

	if len(keys) < 1 {
		return errors.New("The query for SSH keys returned no results. Please modify the search criteria and try again")
	}

	if len(keys) > 1 {
		return fmt.Errorf("The query for SSH keys returned %d results. Please make the search criteria more specific and try again", len(keys))
	}

	d.SetId(keys[0].ID)
	d.Set("name", keys[0].Name)
	d.Set("public_key", keys[0].Key)
	return nil
}
