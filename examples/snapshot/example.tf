// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the OS ID for snapshots.
data "vultr_os" "snapshot" {
  filter {
    name   = "family"
    values = ["snapshot"]
  }
}

// Find the snapshot ID for a Kubernetes master.
data "vultr_snapshot" "master" {
  description_regex = "master"
}

// Find the ID of the Silicon Valley region.
data "vultr_region" "silicon_valley" {
  filter {
    name   = "name"
    values = ["Silicon Valley"]
  }
}

// Find the ID for a starter plan.
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
resource "vultr_instance" "master" {
  name        = "master"
  hostname    = "master"
  region_id   = "${data.vultr_region.silicon_valley.id}"
  plan_id     = "${data.vultr_plan.starter.id}"
  os_id       = "${data.vultr_os.snapshot.id}"
  snapshot_id = "${data.vultr_snapshot.master.id}"
}
