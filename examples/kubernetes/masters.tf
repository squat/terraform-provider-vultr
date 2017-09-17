// Configure the Vultr provider. 
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
#provider "vultr" {
#api_key = "<your-vultr-api-key>"
#}

// Find the ID of the New Jersey region.
data "vultr_region" "new_jersey" {
  filter {
    name   = "state"
    values = ["NJ"]
  }
}

// Find the ID for installing a custom ISO.
data "vultr_os" "custom" {
  filter {
    name   = "family"
    values = ["iso"]
  }
}

// Find the ID for a plan.
data "vultr_plan" "starter" {
  filter {
    name   = "price_per_month"
    values = ["10.00"]
  }

  filter {
    name   = "ram"
    values = ["2048"]
  }
}

// Create a Vultr virtual machine.
resource "vultr_instance" "masters" {
  count              = "${var.master_count}"
  name               = "${var.cluster_name}-master-${count.index}"
  hostname           = "${var.cluster_name}-master-${count.index}"
  region_id          = "${data.vultr_region.new_jersey.id}"
  plan_id            = "${data.vultr_plan.starter.id}"
  os_id              = "${data.vultr_os.custom.id}"
  name               = "${var.cluster_name}-master-${count.index}"
  tag                = "container-linux"
  firewall_group_id  = "${vultr_firewall_group.cluster.id}"
  user_data          = "${data.ct_config.master_ipxe_ignition.rendered}"
  startup_script_id  = "${vultr_startup_script.ipxe.id}"
  private_networking = true
}

resource "vultr_startup_script" "ipxe" {
  type    = "pxe"
  name    = "${var.cluster_name}"
  content = "${file("${path.module}/resources/ipxe")}"
}

// Create a new Vultr DNS domain for the cluster.
resource "vultr_dns_domain" "api" {
  domain = "${var.k8s_domain_name}"
  ip     = "${element(vultr_instance.masters.*.ipv4_address, count.index)}"
}

// Create a DNS record for the API.
resource "vultr_dns_record" "api" {
  count  = "${var.master_count}"
  domain = "${var.k8s_domain_name}"
  name   = "${var.cluster_name}-api"
  type   = "A"
  data   = "${element(vultr_instance.masters.*.ipv4_address, count.index)}"
  ttl    = 300
}
// Output all of the virtual machine's IPv4 addresses to STDOUT when the infrastructure is ready.
output ip_addresses {
  value = "${vultr_instance.masters.*.ipv4_address}"
}
