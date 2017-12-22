module "bootkube" {
  source = "git://https://github.com/poseidon/terraform-render-bootkube.git?ref=v0.9.1"

  cluster_name                  = "${var.cluster_name}"
  api_servers                   = ["${var.cluster_name}-api.${var.k8s_domain_name}"]
  etcd_servers                  = ["http://127.0.0.1:2379"]
  asset_dir                     = "${var.asset_dir}"
  pod_cidr                      = "${var.pod_cidr}"
  service_cidr                  = "${var.service_cidr}"
  experimental_self_hosted_etcd = true
}
