// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the OS ID for applications.
data "vultr_os" "application" {
  filter {
    name   = "family"
    values = ["application"]
  }
}

// Find the application ID for OpenVPN.
data "vultr_application" "openvpn" {
  filter {
    name   = "short_name"
    values = ["openvpn"]
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

// Create a Vultr virtual machine.
resource "vultr_instance" "openvpn" {
  name           = "openvpn"
  hostname       = "openvpn"
  region_id      = data.vultr_region.silicon_valley.id
  plan_id        = data.vultr_plan.starter.id
  os_id          = data.vultr_os.application.id
  application_id = data.vultr_application.openvpn.id
}
