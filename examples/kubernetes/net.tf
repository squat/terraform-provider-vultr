// Create a new firewall group.
resource "vultr_firewall_group" "cluster" {
  description = "${var.cluster_name}"
}

// Add a firewall rule to the group allowing SSH access.
resource "vultr_firewall_rule" "ssh" {
  firewall_group_id = "${vultr_firewall_group.cluster.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}

// Add a firewall rule to the group allowing HTTPS access.
resource "vultr_firewall_rule" "https" {
  firewall_group_id = "${vultr_firewall_group.cluster.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
}

// Add a firewall rule to the group allowing HTTP access.
resource "vultr_firewall_rule" "http" {
  firewall_group_id = "${vultr_firewall_group.cluster.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80
}

// Add a firewall rule to the group allowing ICMP.
resource "vultr_firewall_rule" "icmp" {
  firewall_group_id = "${vultr_firewall_group.cluster.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "icmp"
}
