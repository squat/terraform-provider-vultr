// Create a Vultr virtual machine.
resource "vultr_instance" "workers" {
  count              = "${var.worker_count}"
  name               = "${var.cluster_name}-worker-${count.index}"
  hostname           = "${var.cluster_name}-worker-${count.index}"
  region_id          = "${data.vultr_region.new_jersey.id}"
  plan_id            = "${data.vultr_plan.starter.id}"
  os_id              = "${data.vultr_os.custom.id}"
  name               = "${var.cluster_name}-worker-${count.index}"
  tag                = "container-linux"
  firewall_group_id  = "${vultr_firewall_group.cluster.id}"
  user_data          = "${data.ct_config.worker_ipxe_ignition.rendered}"
  startup_script_id  = "${vultr_startup_script.ipxe.id}"
  private_networking = true
}

// Create a DNS record for the workers for ingress.
resource "vultr_dns_record" "workers" {
  count  = "${var.worker_count}"
  domain = "${var.k8s_domain_name}"
  name   = "${var.cluster_name}-workers"
  type   = "A"
  data   = "${element(vultr_instance.workers.*.ipv4_address, count.index)}"
  ttl    = 300
}
