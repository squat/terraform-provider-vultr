package vultr

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageCreate,
		Read:   resourceBlockStorageRead,
		Update: resourceBlockStorageUpdate,
		Delete: resourceBlockStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cost_per_month": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	name := d.Get("name").(string)
	regionID := d.Get("region_id").(int)
	size := d.Get("size").(int)

	log.Printf("[INFO] Creating new block storage")
	storage, err := client.CreateBlockStorage(name, regionID, size)
	if err != nil {
		return fmt.Errorf("Error creating block storage: %v", err)
	}
	d.SetId(storage.ID)

	instance := d.Get("instance")
	if instance != "" {
		if err := client.AttachBlockStorage(d.Id(), instance.(string)); err != nil {
			return fmt.Errorf("Error attaching newly created block storage (%s) to instance %q: %v", d.Id(), instance.(string), err)
		}
	}

	return resourceBlockStorageRead(d, meta)
}

func resourceBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	storage, err := client.GetBlockStorage(d.Id())
	if err != nil {
		if strings.HasPrefix(err.Error(), "Invalid block storage") {
			log.Printf("[WARN] Removing block storage (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting block storage (%s): %v", d.Id(), err)
	}

	d.Set("cost_per_month", storage.Cost)
	d.Set("size", storage.SizeGB)
	d.Set("name", storage.Name)
	d.Set("region_id", storage.RegionID)
	d.Set("status", storage.Status)
	d.Set("instance", storage.AttachedTo)

	return nil
}

func resourceBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	d.Partial(true)

	if d.HasChange("name") {
		log.Printf("[INFO] Renaming block storage (%s)", d.Id())
		_, new := d.GetChange("name")
		if err := client.LabelBlockStorage(d.Id(), new.(string)); err != nil {
			return fmt.Errorf("Error renaming block storage (%s) to %q: %v", d.Id(), new.(string), err)
		}
		d.SetPartial("name")
	}

	if d.HasChange("size") {
		log.Printf("[INFO] Resizing block storage (%s)", d.Id())
		_, new := d.GetChange("size")
		if err := client.ResizeBlockStorage(d.Id(), new.(int)); err != nil {
			return fmt.Errorf("Error resizing block storage (%s) to %q: %v", d.Id(), new.(int), err)
		}
		d.SetPartial("size")
	}

	if d.HasChange("instance") {
		old, new := d.GetChange("instance")
		if old != "" {
			log.Printf("[INFO] Detaching block storage (%s)", d.Id())
			if err := client.DetachBlockStorage(d.Id()); err != nil {
				return fmt.Errorf("Error detaching block storage (%s): %v", d.Id(), err)
			}
		}
		if new != "" {
			log.Printf("[INFO] Attaching block storage (%s)", d.Id())
			if err := client.AttachBlockStorage(d.Id(), new.(string)); err != nil {
				return fmt.Errorf("Error attaching block storage (%s) to %q: %v", d.Id(), new.(string), err)
			}
		}
		d.SetPartial("instance")
	}

	return resourceBlockStorageRead(d, meta)
}

func resourceBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying block storage (%s)", d.Id())

	instance := d.Get("instance").(string)
	if instance != "" {
		// We need to detach block storage before deleting it
		log.Printf("[INFO] Dettaching block storage (%s) before deleting it.", d.Id())
		if err := client.DetachBlockStorage(d.Id()); err != nil {
			return fmt.Errorf("Error detaching block storage (%s): %v", d.Id(), err)
		}
	}
	if err := client.DeleteBlockStorage(d.Id()); err != nil {
		return fmt.Errorf("Error destroying block storage (%s): %v", d.Id(), err)
	}

	return nil
}
