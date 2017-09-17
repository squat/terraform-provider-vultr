data "ct_config" "master_ipxe_ignition" {
  pretty_print = false
  content      = "${data.template_file.master_ipxe_ignition.rendered}"
}

data "template_file" "master_ipxe_ignition" {
  template = "${file("${path.module}/resources/ipxe_ignition.yaml")}"

  vars {
    ignition = "${data.ct_config.master_ignition.rendered}"
  }
}

data "ct_config" "master_ignition" {
  pretty_print = false
  content      = "${data.template_file.master_config.rendered}"
}

# Master Container Linux Config.
data "template_file" "master_config" {
  template = "${file("${path.module}/resources/master.yaml")}"

  vars {
    k8s_dns_service_ip      = "${cidrhost(var.service_cidr, 10)}"
    kubeconfig_ca_cert      = "${module.bootkube.ca_cert}"
    kubeconfig_kubelet_cert = "${module.bootkube.kubelet_cert}"
    kubeconfig_kubelet_key  = "${module.bootkube.kubelet_key}"
    kubeconfig_server       = "${module.bootkube.server}"
    ssh_authorized_key      = "${file("/home/squat/lserven.ssh")}"
  }
}

data "ct_config" "worker_ipxe_ignition" {
  pretty_print = false
  content      = "${data.template_file.worker_ipxe_ignition.rendered}"
}

data "template_file" "worker_ipxe_ignition" {
  template = "${file("${path.module}/resources/ipxe_ignition.yaml")}"

  vars {
    ignition = "${data.ct_config.worker_ignition.rendered}"
  }
}

data "ct_config" "worker_ignition" {
  pretty_print = false
  content      = "${data.template_file.worker_config.rendered}"
}

# Worker Container Linux Config.
data "template_file" "worker_config" {
  template = "${file("${path.module}/resources/worker.yaml")}"

  vars = {
    k8s_dns_service_ip      = "${cidrhost(var.service_cidr, 10)}"
    kubeconfig_ca_cert      = "${module.bootkube.ca_cert}"
    kubeconfig_kubelet_cert = "${module.bootkube.kubelet_cert}"
    kubeconfig_kubelet_key  = "${module.bootkube.kubelet_key}"
    kubeconfig_server       = "${module.bootkube.server}"
    ssh_authorized_key      = "${file("/home/squat/lserven.ssh")}"
  }
}
