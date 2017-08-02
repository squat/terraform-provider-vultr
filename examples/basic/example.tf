// Configure the Vultr provider. 
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
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
  name              = "basic"
  region_id         = 12                                    // Silicon Valley
  plan_id           = "${data.vultr_plan.starter.id}"       // $5
  os_id             = "${data.vultr_os.container_linux.id}" // CoreOS Container Linux stable
  ssh_keys          = ["${vultr_ssh_key.example.name}"]
  hostname          = "basic"
  tag               = "container-linux"
  firewall_group_id = "${vultr_firewall_group.example.id}"
}

// Create a new firewall group.
resource "vultr_firewall_group" "example" {
  description = "example group"
}

// Add a firewall rule to the group allowing SSH access.
resource "vultr_firewall_rule" "ssh" {
  firewall_group_id = "${vultr_firewall_group.example.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}

// Add a firewall rule to the group allowing ICMP.
resource "vultr_firewall_rule" "icmp" {
  firewall_group_id = "${vultr_firewall_group.example.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "icmp"
}

// Create a new SSH key.
resource "vultr_ssh_key" "example" {
  name       = "squat"
  public_key = "${file("~/lserven.ssh")}"
}

// Add two extra IPv4 addresses to the virtual machine.
resource "vultr_ipv4" "example" {
  instance_id = "${vultr_instance.example.id}"
  reboot      = false
  count       = 2
}

// Output all of the virtual machine's IPv4 addresses to STDOUT when the infrastructure is ready.
output ip_addresses {
  value = "${concat(vultr_ipv4.example.*.ipv4_address, list(vultr_instance.example.ipv4_address))}"
}
