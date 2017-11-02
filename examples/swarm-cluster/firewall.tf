
// Create a firewall.


resource "vultr_firewall_group" "swarm_manager" {
  description = "swarm_manager"
}

resource "vultr_firewall_rule" "ssh_accept_manager" {
  firewall_group_id = "${vultr_firewall_group.swarm_manager.id}"

  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}

resource "vultr_firewall_rule" "http_accept" {
  firewall_group_id = "${vultr_firewall_group.swarm_manager.id}"

  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80

}

resource "vultr_firewall_rule" "https_accept" {
  firewall_group_id = "${vultr_firewall_group.swarm_manager.id}"

  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443

}

resource "vultr_firewall_rule" "ovpn_accept" {
  firewall_group_id = "${vultr_firewall_group.swarm_manager.id}"

  cidr_block        = "0.0.0.0/0"
  protocol          = "udp"
  from_port         = 1194
  to_port           = 1194
}





resource "vultr_firewall_group" "swarm_worker" {
  description = "swarm_worker"
}

resource "vultr_firewall_rule" "ssh_accept_workers" {
  firewall_group_id = "${vultr_firewall_group.swarm_worker.id}"

  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}
