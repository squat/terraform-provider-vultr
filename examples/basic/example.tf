// Configure the Vultr provider. 
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Create a Vultr virtual machine.
resource "vultr_instance" "example" {
  name              = "basic"
  region_id         = 12                                   // Silicon Valley
  plan_id           = 201                                  // $5
  os_id             = 179                                  // CoreOS Container Linux stable
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
