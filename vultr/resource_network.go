package vultr

import (
	"fmt"
	"log"
	"net"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Delete: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var cidrBlock *net.IPNet
	var err error
	if _, ok := d.GetOkExists("cidr_block"); ok {
		_, cidrBlock, err = net.ParseCIDR(d.Get("cidr_block").(string))
		if err != nil {
			return fmt.Errorf("Error parsing %q for netword: %v", "cidr_block", err)
		}
	}
	description := d.Get("description").(string)
	regionID := d.Get("region_id").(int)

	log.Printf("[INFO] Creating new network")
	network, err := client.CreateNetwork(regionID, description, cidrBlock)
	if err != nil {
		return fmt.Errorf("Error creating network: %v", err)
	}

	d.SetId(network.ID)

	return resourceNetworkRead(d, meta)
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	networks, err := client.GetNetworks()
	if err != nil {
		return fmt.Errorf("Error getting network (%s): %v", d.Id(), err)
	}

	var network *lib.Network
	for _, n := range networks {
		if n.ID == d.Id() {
			network = &n
			break
		}
	}

	if network == nil {
		log.Printf("[WARN] Removing network (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("cidr_block", networkToCIDR(network))
	d.Set("description", network.Description)
	d.Set("region_id", network.RegionID)

	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying network (%s)", d.Id())

	if err := client.DeleteNetwork(d.Id()); err != nil {
		return fmt.Errorf("Error destroying network (%s): %v", d.Id(), err)
	}

	return nil
}

func networkToCIDR(net *lib.Network) string {
	return fmt.Sprintf("%s/%d", net.V4Subnet, net.V4SubnetMask)
}
