package vultr

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	osIDSnapshot = 164
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"auto_backups": {
				Type:     schema.TypeBool,
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

			"ipv4_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_mask": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_private_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"ipv6_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"networks": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"notify_activate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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

			"server_state": {
				Type:     schema.TypeString,
				Computed: true,
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
	_, appOK := d.GetOk("application_id")
	_, osOK := d.GetOk("os_id")
	_, snapshotOK := d.GetOk("snapshot_id")
	// At most one.
	if appOK && snapshotOK {
		return fmt.Errorf("Only one of %q and %q may be provided but not both", "application_id", "snapshot_id")
	}
	// Exactly one.
	if osOK == snapshotOK {
		return fmt.Errorf("One of %q and %q must be provided but not both", "os_id", "snapshot_id")
	}

	client := meta.(*Client)
	options := &lib.ServerOptions{
		AppID:                d.Get("application_id").(string),
		AutoBackups:          d.Get("auto_backups").(bool),
		DontNotifyOnActivate: !d.Get("notify_activate").(bool),
		FirewallGroupID:      d.Get("firewall_group_id").(string),
		Hostname:             d.Get("hostname").(string),
		IPV6:                 d.Get("ipv6").(bool),
		PrivateNetworking:    d.Get("private_networking").(bool),
		Script:               d.Get("startup_script_id").(int),
		Snapshot:             d.Get("snapshot_id").(string),
		Tag:                  d.Get("tag").(string),
		UserData:             d.Get("user_data").(string),
	}

	name := d.Get("name").(string)
	var osID int
	if snapshotOK {
		osID = osIDSnapshot
	} else {
		osID = d.Get("os_id").(int)
	}
	planID := d.Get("plan_id").(int)
	regionID := d.Get("region_id").(int)

	netIDs := make([]string, d.Get("network_ids.#").(int))
	for i, id := range d.Get("network_ids").([]interface{}) {
		netIDs[i] = id.(string)
	}
	options.Networks = netIDs

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
		if strings.HasPrefix(err.Error(), "Invalid server") {
			log.Printf("[WARN] Removing instance (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting instance (%s): %v", d.Id(), err)
	}

	networks, err := client.ListPrivateNetworksForServer(d.Id())
	if err != nil {
		return fmt.Errorf("Error getting private networks for instance (%s): %v", d.Id(), err)
	}
	nets := make(map[string]string)
	var networkIDs []string
	for _, n := range networks {
		nets[n.ID] = n.IPAddress
		networkIDs = append(networkIDs, n.ID)
	}

	osID, err := strconv.Atoi(instance.OSID)
	if err != nil {
		return fmt.Errorf("OS ID must be an integer: %v", err)
	}

	var size int
	if instance.InternalIP != "" {
		ipv4s, err := client.ListIPv4(d.Id())
		if err != nil {
			return fmt.Errorf("Error getting IPv4 networks for instance (%s): %v", d.Id(), err)
		}
		err = func() error {
			for _, n := range ipv4s {
				if n.IP == instance.InternalIP {
					size, _ = parseIPv4Mask(n.Netmask).Size()
					return nil
				}
			}
			return fmt.Errorf("no matching address for %q in IPv4 list", instance.InternalIP)
		}()
		if err != nil {
			return fmt.Errorf("Error finding private IPv4 subnet mask for instance (%s): %v", d.Id(), err)
		}
	}

	d.Set("application_id", instance.AppID)
	d.Set("auto_backups", instance.AutoBackups)
	d.Set("cost_per_month", instance.Cost)
	d.Set("default_password", instance.DefaultPassword)
	d.Set("disk", instance.Disk)
	d.Set("firewall_group_id", instance.FirewallGroupID)
	d.Set("ipv4_address", instance.MainIP)
	d.Set("ipv4_gateway", instance.GatewayV4)
	d.Set("ipv4_mask", instance.NetmaskV4)
	d.Set("ipv4_private_cidr", fmt.Sprintf("%s/%d", instance.InternalIP, size))
	d.Set("name", instance.Name)
	d.Set("networks", nets)
	d.Set("network_ids", networkIDs)
	d.Set("os_id", osID)
	d.Set("plan_id", instance.PlanID)
	d.Set("power_status", instance.PowerStatus)
	d.Set("ram", instance.RAM)
	d.Set("region_id", instance.RegionID)
	d.Set("status", instance.Status)
	d.Set("server_state", instance.ServerState)
	d.Set("tag", instance.Tag)
	d.Set("vcpus", instance.VCpus)

	var ipv6s []string
	for _, net := range instance.V6Networks {
		ipv6s = append(ipv6s, net.MainIP)
	}
	d.Set("ipv6_addresses", ipv6s)

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

	if d.HasChange("application_id") {
		log.Printf("[INFO] Updating instance (%s) application", d.Id())
		old, new := d.GetChange("application_id")
		if err := changeApplication(d.Id(), "instance", new.(string), client.ChangeApplicationofServer, client.ListApplicationsforServer); err != nil {
			return err
		}
		if _, err := waitForResourceState(d, meta, "instance", "application_id", resourceInstanceRead, new.(string), []string{"", old.(string)}); err != nil {
			return err
		}
		d.SetPartial("application_id")
	}

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

	if d.HasChange("network_ids") {
		log.Printf("[INFO] Updating instance (%s) networks", d.Id())
		old, new := d.GetChange("network_ids")
		oldIDs := make([]string, len(old.([]interface{})))
		for i, id := range old.([]interface{}) {
			oldIDs[i] = id.(string)
		}
		newIDs := make([]string, len(new.([]interface{})))
		for i, id := range new.([]interface{}) {
			newIDs[i] = id.(string)
		}
		add := stringsDiff(oldIDs, newIDs)
		del := stringsDiff(newIDs, oldIDs)
		for _, n := range add {
			if err := client.EnablePrivateNetworkForServer(d.Id(), n); err != nil {
				return fmt.Errorf("Error attaching instance (%s) to private network %q: %v", d.Id(), n, err)
			}
		}
		for _, n := range del {
			if err := client.DisablePrivateNetworkForServer(d.Id(), n); err != nil {
				return fmt.Errorf("Error detaching instance (%s) to private network %q: %v", d.Id(), n, err)
			}
		}
		d.SetPartial("network_ids")
	}

	if d.HasChange("os_id") {
		log.Printf("[INFO] Updating instance (%s) OS", d.Id())
		old, new := d.GetChange("os_id")
		if err := changeOS(d.Id(), "instance", new.(int), client.ChangeOSofServer, client.ListOSforServer); err != nil {
			return err
		}
		if _, err := waitForResourceState(d, meta, "instance", "os_id", resourceInstanceRead, strconv.FormatInt(int64(new.(int)), 10), []string{"", strconv.FormatInt(int64(old.(int)), 10)}); err != nil {
			return err
		}
		d.SetPartial("os_id")
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

	for {
		// The Vultr API does not allow us to directly check the a server's state if it is being destroyed.
		// However, we can infer this state from the responses of the destroy endpoint.
		// If there was no error, then we need to try destroying again.
		if err := client.DeleteServer(d.Id()); err != nil {
			// Server is locked and still being destroyed. We need to try again.
			if err.Error() == "Unable to destroy server: Unable to remove VM: Server is currently locked" {
				continue
			}

			// Server is pending destruction. We need to wait and try again.
			if err.Error() == "Unable to destroy server: Server is already pending destruction." {
				continue
			}

			// Server does not exist so it has been deleted.
			if strings.HasPrefix(err.Error(), "Invalid server") {
				break
			}
			// There was a legitimate error.
			return fmt.Errorf("Error destroying instance (%s): %v", d.Id(), err)
		}
	}

	return nil
}

// changeOS will try to change the OS of a instance or bare metal instance.
// If there is an error, it will return an error with the list of valid OSs.
func changeOS(id, resourceType string, new int, change func(string, int) error, list func(string) ([]lib.OS, error)) error {
	if err := change(id, new); err != nil {
		var validOS string
		os, oserr := list(id)
		if oserr != nil {
			log.Printf("[Error] failed to get available OSs for %s (%s)", resourceType, id)
		} else {
			var oss []string
			for i := range os {
				oss = append(oss, strconv.FormatInt(int64(os[i].ID), 10))
			}
			validOS = fmt.Sprintf(" Valid OSs are %s", strings.Join(oss, ", "))
		}
		return fmt.Errorf("Error changing OS of %s (%s) to %d: %v%s", resourceType, id, new, err, validOS)
	}
	return nil
}

// changeApplication will try to change the appliation of an instance or bare metal instance.
// If there is an error, it will return an error with the list of valid applications.
func changeApplication(id, resourceType string, new string, change func(string, string) error, list func(string) ([]lib.Application, error)) error {
	if err := change(id, new); err != nil {
		var validApp string
		app, apperr := list(id)
		if apperr != nil {
			log.Printf("[Error] failed to get available applications for %s (%s)", resourceType, id)
		} else {
			var apps []string
			for i := range app {
				apps = append(apps, app[i].ID)
			}
			validApp = fmt.Sprintf(" Valid applications are %s", strings.Join(apps, ", "))
		}
		return fmt.Errorf("Error changing application of %s (%s) to %s: %v%s", resourceType, id, new, err, validApp)
	}
	return nil
}

func parseIPv4Mask(s string) net.IPMask {
	mask := net.ParseIP(s)
	if mask == nil {
		return nil
	}
	return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
}

// stringsDiff returns the strings that are in after and not
// in before.
func stringsDiff(before []string, after []string) []string {
	var diff []string
	b := map[string]struct{}{}
	for i := range before {
		b[before[i]] = struct{}{}
	}
	for i := range after {
		if _, ok := b[after[i]]; !ok {
			diff = append(diff, after[i])
		}
	}
	return diff
}
