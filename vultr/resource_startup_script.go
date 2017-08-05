package vultr

import (
	"fmt"
	"log"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStartupScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceStartupScriptCreate,
		Read:   resourceStartupScriptRead,
		Update: resourceStartupScriptUpdate,
		Delete: resourceStartupScriptDelete,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStartupScriptType,
			},
		},
	}
}

func resourceStartupScriptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	content := d.Get("content").(string)
	name := d.Get("name").(string)
	var scriptType string
	if _, typeOk := d.GetOk("type"); typeOk {
		scriptType = d.Get("type").(string)
	} else {
		scriptType = "boot"
	}

	log.Printf("[INFO] Creating new startup script")
	script, err := client.CreateStartupScript(name, content, scriptType)
	if err != nil {
		return fmt.Errorf("Error creating startup script: %v", err)
	}

	d.SetId(script.ID)

	return resourceStartupScriptRead(d, meta)
}

func resourceStartupScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	script, err := client.GetStartupScript(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting startup script (%s): %v", d.Id(), err)
	}

	d.Set("content", script.Content)
	d.Set("name", script.Name)
	d.Set("type", script.Type)

	return nil
}

func resourceStartupScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	script := lib.StartupScript{
		Content: d.Get("content").(string),
		ID:      d.Id(),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
	}

	if err := client.UpdateStartupScript(script); err != nil {
		return fmt.Errorf("Error updating startup script (%s): %v", d.Id(), err)
	}

	return resourceStartupScriptRead(d, meta)
}

func resourceStartupScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Destroying startup script (%s)", d.Id())

	if err := client.DeleteStartupScript(d.Id()); err != nil {
		return fmt.Errorf("Error destroying startup script (%s): %v", d.Id(), err)
	}

	return nil
}
