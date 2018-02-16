package vultr

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cost_per_month": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"firewall_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_private_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"os_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"plan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"power_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_networking": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"startup_script_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	options := &lib.ServerOptions{
		AppID:             d.Get("application_id").(string),
		FirewallGroupID:   d.Get("firewall_group_id").(string),
		Hostname:          d.Get("hostname").(string),
		IPV6:              d.Get("ipv6").(bool),
		PrivateNetworking: d.Get("private_networking").(bool),
		Script:            d.Get("startup_script_id").(int),
		Snapshot:          d.Get("snapshot_id").(string),
		Tag:               d.Get("tag").(string),
		UserData:          d.Get("user_data").(string),
	}

	name := d.Get("name").(string)
	osID := d.Get("os_id").(int)
	planID := d.Get("plan_id").(int)
	regionID := d.Get("region_id").(int)

	keyIDs := make([]string, d.Get("ssh_key_ids.#").(int))
	for i, id := range d.Get("ssh_key_ids").([]interface{}) {
		keyIDs[i] = id.(string)
	}
	options.SSHKey = strings.Join(keyIDs, ",")

	log.Printf("[INFO] Creating new instance")
	instance, err := client.CreateServer(name, regionID, planID, osID, options)
	if err != nil {
		return fmt.Errorf("Error creating instance: %v", err)
	}
	d.SetId(instance.ID)

	if _, err := waitForResourceState(d, meta, "instance", "status", resourceInstanceRead, "active", []string{"pending"}); err != nil {
		return err
	}
	if _, err := waitForResourceState(d, meta, "instance", "power_status", resourceInstanceRead, "running", []string{"starting", "stopped"}); err != nil {
		return err
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	instance, err := client.GetServer(d.Id())
	if err != nil {
		if err.Error() == "Invalid server." {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting instance (%s): %v", d.Id(), err)
	}

	d.Set("cost_per_month", instance.Cost)
	d.Set("default_password", instance.DefaultPassword)
	d.Set("disk", instance.Disk)
	if instance.FirewallGroupID == "" {
		d.Set("firewall_group_id", 0)
	} else {
		d.Set("firewall_group_id", instance.FirewallGroupID)
	}
	d.Set("ipv4_address", instance.MainIP)
	d.Set("ipv4_private_address", instance.InternalIP)
	d.Set("name", instance.Name)
	d.Set("os_id", instance.OSID)
	d.Set("plan_id", instance.PlanID)
	d.Set("power_status", instance.PowerStatus)
	d.Set("ram", instance.RAM)
	d.Set("region_id", instance.RegionID)
	d.Set("status", instance.Status)
	d.Set("tag", instance.Tag)
	d.Set("vcpus", instance.VCpus)

	// Initialize the connection information.
	d.SetConnInfo(map[string]string{
		"host":     instance.MainIP,
		"password": instance.DefaultPassword,
		"type":     "ssh",
		"user":     "root",
	})

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	d.Partial(true)

	if d.HasChange("firewall_group_id") {
		log.Printf("[INFO] Updating instance (%s) firewall group", d.Id())
		old, new := d.GetChange("firewall_group_id")
		if err := client.SetFirewallGroup(d.Id(), new.(string)); err != nil {
			return fmt.Errorf("Error changing instance (%s) firewall group to %q: %v", d.Id(), new.(string), err)
		}
		if _, err := waitForResourceState(d, meta, "instance", "firewall_group_id", resourceInstanceRead, new.(string), []string{old.(string)}); err != nil {
			return err
		}
	}

	if d.HasChange("name") {
		log.Printf("[INFO] Updating instance (%s) name", d.Id())
		old, new := d.GetChange("name")
		if err := client.RenameServer(d.Id(), new.(string)); err != nil {
			return fmt.Errorf("Error renaming instance (%s) to %q: %v", d.Id(), new.(string), err)
		}
		if _, err := waitForResourceState(d, meta, "instance", "name", resourceInstanceRead, new.(string), []string{"", old.(string)}); err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("tag") {
		log.Printf("[INFO] Updating instance (%s) tag", d.Id())
		old, new := d.GetChange("tag")
		if err := client.TagServer(d.Id(), new.(string)); err != nil {
			return fmt.Errorf("Error tagging instance (%s) with %q: %v", d.Id(), new.(string), err)
		}
		if _, err := waitForResourceState(d, meta, "instance", "tag", resourceInstanceRead, new.(string), []string{"", old.(string)}); err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	d.Partial(false)

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying instance (%s)", d.Id())

	if err := client.DeleteServer(d.Id()); err != nil {
		return fmt.Errorf("Error destroying instance (%s): %v", d.Id(), err)
	}

	return nil
}
