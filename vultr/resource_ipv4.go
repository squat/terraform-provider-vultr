package vultr

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIPV4() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPV4Create,
		Read:   resourceIPV4Read,
		Update: resourceIPV4Update,
		Delete: resourceIPV4Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"reboot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"reverse_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIPV4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	instance := d.Get("instance_id").(string)
	reboot := d.Get("reboot").(bool)

	log.Printf("[INFO] Creating new IPv4 address")

	vultrMutexKV.Lock(instance)
	ips, err := client.ListIPv4(instance)
	if err != nil {
		return fmt.Errorf("Error listing IPv4 addresses: %v", err)
	}
	if err := client.CreateIPv4(instance, reboot); err != nil {
		return fmt.Errorf("Error creating IPv4 address: %v", err)
	}
	newIPs, err := client.ListIPv4(instance)
	if err != nil {
		return fmt.Errorf("Error re-listing IPv4 addresses: %v", err)
	}
	vultrMutexKV.Unlock(instance)

	var ip *lib.IPv4
	for i := range newIPs {
		var found bool
		for j := range ips {
			if newIPs[i] == ips[j] {
				found = true
				break
			}
		}
		if !found {
			ip = &newIPs[i]
			break
		}
	}
	if ip == nil {
		return errors.New("Error finding created IPv4 address")
	}

	d.SetId(fmt.Sprintf("%s/%s", instance, ip.IP))

	return resourceIPV4Read(d, meta)
}

func resourceIPV4Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	instance, id, err := idFromData(d)
	if err != nil {
		return err
	}

	ips, err := client.ListIPv4(instance)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Invalid server.") {
			log.Printf("[WARN] Removing IPv4 address (%s) because the attached instance (%s) is gone", d.Id(), instance)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting IPv4 addresses: %v", err)
	}
	var ip *lib.IPv4
	for i := range ips {
		if ips[i].IP == id {
			ip = &ips[i]
			break
		}
	}
	if ip == nil {
		log.Printf("[WARN] Removing IPv4 address (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("gateway", ip.Gateway)
	d.Set("ipv4_address", ip.IP)
	d.Set("netmask", ip.Netmask)
	d.Set("reverse_dns", ip.ReverseDNS)

	return nil
}

func resourceIPV4Update(d *schema.ResourceData, meta interface{}) error {
	return resourceIPV4Read(d, meta)
}

func resourceIPV4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	instance, id, err := idFromData(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Destroying IPv4 address (%s)", d.Id())

	if err := client.DeleteIPv4(instance, id); err != nil {
		return fmt.Errorf("Error destroying IPv4 address (%s): %v", d.Id(), err)
	}

	return nil
}

// idFromData returns the IPv4 id components from the ResourceData ID,
// which are, in order: the IPv4's associated instance, the actual IP, and any error.
func idFromData(d *schema.ResourceData) (string, string, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return "", "", errors.New("Error parsing IPv4 ID: ID should be of form <instance-id>/<ip-address>")
	}
	return parts[0], parts[1], nil
}
