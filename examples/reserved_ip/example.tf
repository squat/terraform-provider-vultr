// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
#provider "vultr" {
#api_key = "<your-vultr-api-key>"
#}

// Find the ID of the Silicon Valley region.
data "vultr_region" "silicon_valley" {
  filter {
    name   = "name"
    values = ["Silicon Valley"]
  }
}

// Find the ID for CoreOS Container Linux.
data "vultr_os" "container_linux" {
  filter {
    name   = "family"
    values = ["coreos"]
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
resource "vultr_instance" "example" {
  name        = "example"
  region_id   = data.vultr_region.silicon_valley.id
  plan_id     = data.vultr_plan.starter.id
  os_id       = data.vultr_os.container_linux.id
  ssh_key_ids = [vultr_ssh_key.squat.id]
}

// Create a new SSH key.
resource "vultr_ssh_key" "squat" {
  name       = "squat"
  public_key = file("~/lserven.ssh")
}

// Create a reserved IP.
resource "vultr_reserved_ip" "example" {
  name        = "example"
  attached_id = vultr_instance.example.id
  region_id   = data.vultr_region.silicon_valley.id
}
