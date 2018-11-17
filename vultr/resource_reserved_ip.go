package vultr

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceReservedIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceReservedIPCreate,
		Read:   resourceReservedIPRead,
		Update: resourceReservedIPUpdate,
		Delete: resourceReservedIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"attached_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:         schema.TypeString,
				Default:      "v4",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateReservedIPType,
			},
		},
	}
}

func resourceReservedIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	aid := d.Get("attached_id").(string)
	name := d.Get("name").(string)
	regionID := d.Get("region_id").(int)
	ipType := d.Get("type").(string)

	log.Printf("[INFO] Creating new reserved ip")
	rid, err := client.CreateReservedIP(regionID, ipType, name)
	if err != nil {
		return fmt.Errorf("Error creating reserved ip: %v", err)
	}

	d.SetId(rid)

	if aid != "" {
		rip, err := client.GetReservedIP(rid)
		if err != nil {
			return fmt.Errorf("Error getting address for reserved ip (%s): %v", d.Id(), err)
		}
		err = client.AttachReservedIP(reservedIPToCIDR(rip), aid)
		if err != nil {
			return fmt.Errorf("Error attaching reserved ip (%s) to instance (%s): %v", d.Id(), aid, err)
		}
	}

	return resourceReservedIPRead(d, meta)
}

func resourceReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rip, err := client.GetReservedIP(d.Id())
	if err != nil {
		if strings.HasPrefix(err.Error(), fmt.Sprintf("IP with ID %v not found", d.Id())) {
			log.Printf("[WARN] Removing reserved ip (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting reserved ip (%s): %v", d.Id(), err)
	}

	d.Set("attached_id", rip.AttachedTo)
	d.Set("cidr", reservedIPToCIDR(rip))
	d.Set("name", rip.Label)
	d.Set("region_id", rip.RegionID)
	d.Set("type", rip.IPType)

	return nil
}

func resourceReservedIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	if d.HasChange("attached_id") {
		ip := d.Get("cidr").(string)
		log.Printf("[INFO] Updating reserved IP (%s) attachment", d.Id())
		old, new := d.GetChange("attached_id")
		if old.(string) != "" {
			if err := client.DetachReservedIP(old.(string), ip); err != nil {
				return fmt.Errorf("Error detaching reserved IP (%s) from instance (%s): %v", d.Id(), old.(string), err)
			}
		}
		if new.(string) != "" {
			if err := client.AttachReservedIP(ip, new.(string)); err != nil {
				return fmt.Errorf("Error attaching reserved IP (%s) to instance (%s): %v", d.Id(), new.(string), err)
			}
		}
	}

	return resourceReservedIPRead(d, meta)
}

func resourceReservedIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying reserved ip (%s)", d.Id())

	aid := d.Get("attached_id").(string)
	if aid != "" {
		ip := d.Get("cidr").(string)
		if err := client.DetachReservedIP(aid, ip); err != nil {
			return fmt.Errorf("Error detaching reserved IP (%s) from instance (%s): %v", d.Id(), aid, err)
		}
	}

	if err := client.DestroyReservedIP(d.Id()); err != nil {
		return fmt.Errorf("Error destroying reserved ip (%s): %v", d.Id(), err)
	}

	return nil
}

func reservedIPToCIDR(rip lib.IP) string {
	return fmt.Sprintf("%s/%d", rip.Subnet, rip.SubnetSize)
}
