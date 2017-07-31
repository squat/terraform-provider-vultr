package vultr

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Update: resourceSSHKeyUpdate,
		Delete: resourceSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					return strings.TrimSpace(v.(string))
				},
			},
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	log.Printf("[INFO] Creating new SSH key")
	key, err := client.CreateSSHKey(name, publicKey)
	if err != nil {
		return fmt.Errorf("Error creating SSH key: %v", err)
	}
	d.SetId(key.ID)

	return resourceSSHKeyRead(d, meta)
}

func resourceSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	keys, err := client.GetSSHKeys()
	if err != nil {
		return fmt.Errorf("Error getting SSH keys: %v", err)
	}
	var key *lib.SSHKey
	for i := range keys {
		if keys[i].ID == d.Id() {
			key = &keys[i]
			break
		}
		if i == len(keys)-1 {
			log.Printf("[WARN] Removing SSH key (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
	}

	d.Set("name", key.Name)
	d.Set("public_key", key.Key)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Updating SSH key (%s)", d.Id())

	key := lib.SSHKey{
		ID:   d.Id(),
		Key:  d.Get("public_key").(string),
		Name: d.Get("name").(string),
	}
	if err := client.UpdateSSHKey(key); err != nil {
		return fmt.Errorf("Error updating SSH key (%s): %v", d.Id(), err)
	}

	return resourceSSHKeyRead(d, meta)
}

func resourceSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying SSH key (%s)", d.Id())

	if err := client.DeleteSSHKey(d.Id()); err != nil {
		return fmt.Errorf("Error destroying SSH key (%s): %v", d.Id(), err)
	}

	return nil
}
