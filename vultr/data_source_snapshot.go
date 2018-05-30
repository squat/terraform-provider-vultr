package vultr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSnapshotRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegex,
			},

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	filters, filtersOk := d.GetOk("filter")
	descriptionRegex, descriptionRegexOk := d.GetOk("description_regex")

	if !filtersOk && !descriptionRegexOk {
		return fmt.Errorf("One of %q and %q must be provided", "filter", "description_regex")
	}

	snapshots, err := client.GetSnapshots()
	if err != nil {
		return fmt.Errorf("Error getting snapshots: %v", err)
	}

	if filtersOk {
		filter := filterFromSet(filters.(*schema.Set))
		var filteredSnapshots []lib.Snapshot
		for _, snapshot := range snapshots {
			m := structToMap(snapshot)
			if filter.F(m) {
				filteredSnapshots = append(filteredSnapshots, snapshot)
			}
		}
		snapshots = filteredSnapshots
	}

	if descriptionRegexOk {
		var filteredSnapshots []lib.Snapshot
		r := regexp.MustCompile(descriptionRegex.(string))
		for _, snapshot := range snapshots {
			if r.MatchString(snapshot.Description) {
				filteredSnapshots = append(filteredSnapshots, snapshot)
			}
		}
		snapshots = filteredSnapshots
	}

	if len(snapshots) < 1 {
		return errors.New("The query for snapshots returned no results. Please modify the search criteria and try again")
	}

	if len(snapshots) > 1 {
		return fmt.Errorf("The query for snapshots returned %d results. Please make the search criteria more specific and try again", len(snapshots))
	}

	osID, err := strconv.Atoi(snapshots[0].OSID)
	if err != nil {
		return fmt.Errorf("OS ID must be an integer: %v", err)
	}

	d.SetId(snapshots[0].ID)
	d.Set("application_id", snapshots[0].AppID)
	d.Set("description", snapshots[0].Description)
	d.Set("created", snapshots[0].Created)
	d.Set("os_id", osID)
	d.Set("size", snapshots[0].Size)
	d.Set("status", snapshots[0].Status)
	return nil
}
