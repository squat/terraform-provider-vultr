package vultr

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSDomainCreate,
		Read:   resourceDNSDomainRead,
		Delete: resourceDNSDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIPAddress,
			},
		},
	}
}

func resourceDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	domain := d.Get("domain").(string)
	ip := d.Get("ip").(string)

	log.Printf("[INFO] Creating new DNS domain")
	err := client.CreateDNSDomain(domain, ip)
	if err != nil {
		return fmt.Errorf("Error creating DNS domain: %v", err)
	}

	d.SetId(domain)

	return resourceDNSDomainRead(d, meta)
}

func resourceDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	dnsDomains, err := client.GetDNSDomains()
	if err != nil {
		return fmt.Errorf("Error getting DNS domains: %v", err)
	}

	for i := range dnsDomains {
		if dnsDomains[i].Domain == d.Id() {
			return nil
		}
	}

	log.Printf("[WARN] Removing DNS domain (%s) because it is gone", d.Id())
	d.SetId("")

	return nil
}

func resourceDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying DNS domain (%s)", d.Id())

	if err := client.DeleteDNSDomain(d.Id()); err != nil {
		return fmt.Errorf("Error destroying DNS domain (%s): %v", d.Id(), err)
	}

	return nil
}
