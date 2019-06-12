package vultr

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"golang.org/x/net/idna"
)

func resourceReverseName() *schema.Resource {
	return &schema.Resource{
		Create: resourceReverseNameCreate,
		Read:   resourceReverseNameRead,
		Update: resourceReverseNameUpdate,
		Delete: resourceReverseNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return net.ParseIP(old).Equal(net.ParseIP(new))
				},
				ValidateFunc: validation.IPRange(),
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					oldASCII, _ := idna.Registration.ToASCII(old)
					newASCII, _ := idna.Registration.ToASCII(new)

					return oldASCII == newASCII
				},
				StateFunc: func(val interface{}) string {
					ascii, _ := idna.Registration.ToASCII(val.(string))
					return ascii
				},
				ValidateFunc: func(val interface{}, key string) (_ []string, errs []error) {
					_, err := idna.Registration.ToASCII(val.(string))
					if err != nil {
						errs = append(errs, fmt.Errorf("Error validating %s: Expected '%s' to be a valid domain name", key, val))
					}
					return
				},
			},
		},
	}
}

func resourceReverseNameCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Get("instance_id").(string)
	ip := net.ParseIP(d.Get("ip").(string))
	name := d.Get("name").(string)

	log.Printf("[INFO] Setting reverse name for IP (%s)", ip)

	setOrUpdateReverseName(client, id, ip, name)

	d.SetId(fmt.Sprintf("%s/%s", id, ip))

	return resourceReverseNameRead(d, meta)
}

func resourceReverseNameRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id, ip, err := splitID(d.Id())
	if err != nil {
		return err
	}

	d.Set("instance_id", id)
	d.Set("ip", ip.String())

	if isIPv4(ip) {
		ipv4s, err := client.ListIPv4(id)
		if err != nil {
			return fmt.Errorf("Error getting a list of IPv4 addresses for the instance (%s): %v", id, err)
		}

		for _, ipv4 := range ipv4s {
			if ip.Equal(net.ParseIP(ipv4.IP)) {
				d.Set("name", ipv4.ReverseDNS)
				return nil
			}
		}

		log.Printf("[INFO] Can't get reverse name of IPv4 (%s) on instance (%s). Assuming that this IPv4 was deleted.", ip, id)
		d.SetId("")
		return nil
	}

	ipv6s, err := client.ListIPv6ReverseDNS(id)
	if err != nil {
		return fmt.Errorf("Error getting a list of IPv6 addresses for the instance (%s): %v", id, err)
	}

	for _, ipv6 := range ipv6s {
		if ip.Equal(net.ParseIP(ipv6.IP)) {
			d.Set("name", ipv6.ReverseDNS)
			return nil
		}
	}

	log.Printf("[INFO] Can't get reverse name of IPv6 (%s) on instance (%s). Assuming that this IPv6 was deleted.", ip, id)
	d.SetId("")
	return nil
}

func resourceReverseNameUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	if d.HasChange("name") {
		id := d.Get("instance_id").(string)
		ip := net.ParseIP(d.Get("ip").(string))
		log.Printf("[INFO] Updating reverse name for IP (%s)", ip)
		_, new := d.GetChange("name")
		setOrUpdateReverseName(client, id, ip, new.(string))
	}

	return resourceReverseNameRead(d, meta)
}

func setOrUpdateReverseName(client *Client, id string, ip net.IP, name string) error {
	if isIPv4(ip) {
		err := client.SetIPv4ReverseDNS(id, ip.String(), name)
		if err != nil {
			return fmt.Errorf("Error setting reverse name on instance (%s) for IPv4 (%s) to '%s': %v", id, ip, name, err)
		}
		return nil
	}

	err := client.SetIPv6ReverseDNS(id, ip.String(), name)
	if err != nil {
		return fmt.Errorf("Error setting reverse name on instance (%s) for IPv6 (%s) to '%s': %v", id, ip, name, err)
	}
	return nil
}

func resourceReverseNameDelete(d *schema.ResourceData, meta interface{}) error {
	id, ip, err := splitID(d.Id())

	log.Printf("[INFO] Destroying reverse name for IP (%s)", ip)

	if isIPv4(ip) {
		// The reverse name can not be "unset", so doing nothing here
		d.SetId("")
		return nil
	}

	client := meta.(*Client)

	err = client.DeleteIPv6ReverseDNS(id, ip.String())
	if err != nil {
		return fmt.Errorf("Error deleting reverse name for IP (%s): %v", ip, err)
	}

	d.SetId("")
	return nil
}

func splitID(resourceID string) (id string, ip net.IP, err error) {
	strs := strings.SplitN(resourceID, "/", 2)

	if len(strs) != 2 {
		err = fmt.Errorf("Error decoding id '%s', the format should be 'instance_id/IP'", resourceID)
		return
	}

	id = strs[0]
	ip = net.ParseIP(strs[1])
	if ip == nil {
		err = fmt.Errorf("Error parsing '%s' as IP", strs[1])
		return
	}

	return
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}
