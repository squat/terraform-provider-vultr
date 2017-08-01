package vultr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOSRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"surcharge": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"windows": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceOSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !filtersOk && !nameRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "name_regex")
	}

	images, err := client.GetOS()
	if err != nil {
		return fmt.Errorf("Error getting images: %v", err)
	}

	filter := filterFromSet(filters.(*schema.Set))
	if filtersOk {
		var filteredImages []lib.OS
		for _, image := range images {
			m := structToMap(image)
			if filter.F(m) {
				filteredImages = append(filteredImages, image)
			}
		}
		images = filteredImages
	}

	if nameRegexOk {
		var filteredImages []lib.OS
		r := regexp.MustCompile(nameRegex.(string))
		for _, image := range images {
			if r.MatchString(image.Name) {
				filteredImages = append(filteredImages, image)
			}
		}
		images = filteredImages
	}

	if len(images) < 1 {
		return errors.New("The query returned for OS returned no results. Please modify the search criteria and try again")
	}

	if len(images) > 1 {
		return fmt.Errorf("The query returned for OS returned %d results. Please make the search criteria more specific and try again", len(images))
	}

	d.SetId(strconv.Itoa(images[0].ID))
	d.Set("arch", images[0].Arch)
	d.Set("family", images[0].Family)
	d.Set("name", images[0].Name)
	d.Set("surcharge", images[0].Surcharge)
	d.Set("windows", images[0].Windows)
	return nil
}
