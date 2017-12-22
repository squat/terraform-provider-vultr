package vultr

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/schema"
)

const dnsRecordIDFormatErrTemplate = "DNS record ID must conform to <domain>/<record-ID>, where <record-ID> is an integer; got %q"

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Delete: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	data := d.Get("data").(string)
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	var priority int
	if recordType == "MX" || recordType == "SRV" {
		if _, priorityOk := d.GetOk("priority"); priorityOk {
			priority = d.Get("priority").(int)
		} else {
			return fmt.Errorf("Records of type %s and %s require a priority", "MX", "SRV")
		}
	}

	var ttl int
	if _, ttlOk := d.GetOk("ttl"); ttlOk {
		ttl = d.Get("ttl").(int)
	}

	log.Printf("[INFO] Creating new DNS record")
	err := client.CreateDNSRecord(domain, name, recordType, data, priority, ttl)
	if err != nil {
		return fmt.Errorf("Error creating DNS record: %v", err)
	}

	records, err := client.GetDNSRecords(domain)
	if err != nil {
		return fmt.Errorf("Error getting DNS record: %v", err)
	}
	for _, r := range records {
		if r.Data == data && r.Name == name && r.Type == recordType {
			d.SetId(fmt.Sprintf("%s/%d", domain, r.RecordID))
			return resourceDNSRecordRead(d, meta)
		}
	}

	return fmt.Errorf("Error finding DNS record: %v", err)
}

func resourceDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return fmt.Errorf(dnsRecordIDFormatErrTemplate, d.Id())
	}
	domain := idParts[0]
	id, err := strconv.Atoi(idParts[1])
	if err != nil {
		return fmt.Errorf(dnsRecordIDFormatErrTemplate, d.Id())
	}

	records, err := client.GetDNSRecords(domain)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Invalid domain.") {
			log.Printf("[WARN] Removing DNS record (%s) because the domain is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting DNS records: %v", err)
	}

	var record *lib.DNSRecord
	for i := range records {
		if records[i].RecordID == id {
			record = &records[i]
			break
		}
	}

	if record == nil {
		log.Printf("[WARN] Removing DNS record (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("data", record.Data)
	d.Set("domain", domain)
	d.Set("name", record.Name)
	d.Set("priority", record.Priority)
	d.Set("ttl", record.TTL)
	d.Set("type", record.Type)

	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idParts := strings.Split(d.Id(), "/")
	domain := idParts[0]
	id, _ := strconv.Atoi(idParts[1])

	record := lib.DNSRecord{
		Data:     d.Get("data").(string),
		RecordID: id,
		Name:     d.Get("name").(string),
		Priority: d.Get("priority").(int),
		TTL:      d.Get("ttl").(int),
		Type:     d.Get("type").(string),
	}

	err := client.UpdateDNSRecord(domain, record)
	if err != nil {
		return fmt.Errorf("Error updating DNS record (%s): %v", d.Id(), err)
	}

	return resourceDNSRecordRead(d, meta)
}

func resourceDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	idParts := strings.Split(d.Id(), "/")
	domain := idParts[0]
	id, _ := strconv.Atoi(idParts[1])

	log.Printf("[INFO] Destroying DNS record (%s)", d.Id())

	if err := client.DeleteDNSRecord(domain, id); err != nil {
		return fmt.Errorf("Error destroying DNS record (%s): %v", d.Id(), err)
	}

	return nil
}
