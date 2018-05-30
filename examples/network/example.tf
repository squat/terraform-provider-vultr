// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the ID of the Frankfurt region.
data "vultr_region" "frankfurt" {
  filter {
    name   = "name"
    values = ["Frankfurt"]
  }
}

// Find the ID for a starter plan.
data "vultr_plan" "starter" {
  filter {
    name   = "price_per_month"
    values = ["5.00"]
  }

  filter {
    name   = "ram"
    values = ["1024"]
  }
}

// Find the OS ID for Ubuntu.
data "vultr_os" "ubuntu" {
  filter {
    name   = "name"
    values = ["Ubuntu 18.04 x64"]
  }
}

// Create a pair of Vultr private networks.
resource "vultr_network" "network" {
  count       = 2
  cidr_block  = "${cidrsubnet("192.168.0.0/23", 1, count.index)}"
  description = "test_${count.index}"
  region_id   = "${data.vultr_region.frankfurt.id}"
}

// Create a Vultr virtual machine.
resource "vultr_instance" "ubuntu" {
  name        = "ubuntu"
  network_ids = ["${vultr_network.network.*.id}"]
  region_id   = "${data.vultr_region.frankfurt.id}"
  plan_id     = "${data.vultr_plan.starter.id}"
  os_id       = "${data.vultr_os.ubuntu.id}"
}
