package vultr

import (
	"fmt"
	"log"
	"strconv"
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
			switch t := state.(type) {
			case bool:
				return struct{}{}, strconv.FormatBool(state.(bool)), nil
			case int:
				return struct{}{}, strconv.FormatInt(int64(state.(int)), 10), nil
			case string:
				return struct{}{}, state.(string), nil
			default:
				return struct{}{}, "", fmt.Errorf("do not know how to interpret %v of type %T", state, t)
			}
		}
		return nil, "", nil
	}
}
