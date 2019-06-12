package vultr

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceReverseName() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceReverseNameRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return net.ParseIP(old).Equal(net.ParseIP(new))
				},
				ValidateFunc: validation.SingleIP(),
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceReverseNameRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Get("instance_id").(string)
	ip := net.ParseIP(d.Get("ip").(string))

	d.SetId(fmt.Sprintf("%s/%s", id, ip))
	d.Set("instance_id", id)
	d.Set("ip", ip.String())

	if isIPv4(ip) {
		ipv4s, err := client.ListIPv4(id)
		if err != nil {
			return fmt.Errorf("Error getting a list of IPv4 addresses for the instance (%s): %v", id, err)
		}
		if len(ipv4s) < 1 {
			return fmt.Errorf("Error getting a list of IPv4 addresses for the instance (%s): Empty list", id)
		}

		for _, ipv4 := range ipv4s {
			if ip.Equal(net.ParseIP(ipv4.IP)) {
				d.Set("name", ipv4.ReverseDNS)
				return nil
			}
		}

		return fmt.Errorf("Error getting the reverse name of the IPv4 '%s' on the instance (%s)", ip, id)
	}

	ipv6s, err := client.ListIPv6ReverseDNS(id)
	if err != nil {
		return fmt.Errorf("Error getting a list of IPv6 addresses for the instance (%s): %v", id, err)
	}
	if len(ipv6s) < 1 {
		return fmt.Errorf("Error getting a list of IPv6 addresses for the instance (%s): Empty list", id)
	}

	for _, ipv6 := range ipv6s {
		if ip.Equal(net.ParseIP(ipv6.IP)) {
			d.Set("name", ipv6.ReverseDNS)
			return nil
		}
	}

	return fmt.Errorf("Error getting the reverse name of the IPv6 '%s' on the instance (%s)", ip, id)
}
