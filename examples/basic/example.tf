resource "vultr_instance" "example" {
  name              = "basic-example"
  region_id         = 12                                   // Silicon Valley
  plan_id           = 201                                  // $5
  os_id             = 179                                  // CoreOS Container Linux stable
  ssh_keys          = ["${vultr_ssh_key.example.name}"]
  hostname          = "basic"
  tag               = "container-linux"
  firewall_group_id = "${vultr_firewall_group.example.id}"
}

resource "vultr_firewall_group" "example" {
  description = "example group"
}

resource "vultr_firewall_rule" "ssh" {
  firewall_group_id = "${vultr_firewall_group.example.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}

resource "vultr_firewall_rule" "icmp" {
  firewall_group_id = "${vultr_firewall_group.example.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "icmp"
}

resource "vultr_ssh_key" "example" {
  name       = "squat"
  public_key = "${file("~/lserven.ssh")}"
}

resource "vultr_ipv4" "example" {
  instance_id = "${vultr_instance.example.id}"
  reboot      = false
  count       = 2
}

output ip_addresses {
  value = "${concat(vultr_ipv4.example.*.ipv4_address, list(vultr_instance.example.ipv4_address))}"
}
