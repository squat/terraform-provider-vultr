// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the OS ID for Container Linux.
data "vultr_os" "container_linux" {
  filter {
    name   = "family"
    values = ["coreos"]
  }
}

// Find the ID of the Silicon Valley region.
data "vultr_region" "silicon_valley" {
  filter {
    name   = "name"
    values = ["Silicon Valley"]
  }
}

// Find the ID for a starter plan.
data "vultr_bare_metal_plan" "eightcpus" {
  filter {
    name   = "cpu_count"
    values = [8]
  }
}

// Create a Vultr virtual machine.
resource "vultr_bare_metal" "example" {
  name      = "example"
  hostname  = "example"
  region_id = data.vultr_region.silicon_valley.id
  plan_id   = data.vultr_bare_metal_plan.eightcpus.id
  os_id     = data.vultr_os.container_linux.id
}
