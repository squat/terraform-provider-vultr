package vultr

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func waitForResourceState(d *schema.ResourceData, meta interface{}, resourceName, prop string, readFunc schema.ReadFunc, target string, pending []string) (interface{}, error) {
	log.Printf("[INFO] Waiting for %s (%s) to update %s to %q", resourceName, d.Id(), prop, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    resourceStateRefreshFunc(d, meta, prop, readFunc),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	state, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for %s (%s) to have %s %q", resourceName, d.Id(), prop, target)
	}
	return state, nil
}

func resourceStateRefreshFunc(d *schema.ResourceData, meta interface{}, prop string, readFunc schema.ReadFunc) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		if err := readFunc(d, meta); err != nil {
			return nil, "", err
		}

		if state, ok := d.GetOk(prop); ok {
			return struct{}{}, state.(string), nil
		}
		return nil, "", nil
	}
}
