package vultr

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VULTR_API_KEY", nil),
				Description: "The key for API operations. You can retrieve this from the 'API' tab of the 'Account' section  of the Vultr console.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vultr_os":      dataSourceOS(),
			"vultr_plan":    dataSourcePlan(),
			"vultr_ssh_key": dataSourceSSHKey(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vultr_firewall_group": resourceFirewallGroup(),
			"vultr_firewall_rule":  resourceFirewallRule(),
			"vultr_instance":       resourceInstance(),
			"vultr_ipv4":           resourceIPV4(),
			"vultr_startup_script": resourceStartupScript(),
			"vultr_ssh_key":        resourceSSHKey(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: d.Get("api_key").(string),
	}
	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var vultrMutexKV = mutexkv.NewMutexKV()
