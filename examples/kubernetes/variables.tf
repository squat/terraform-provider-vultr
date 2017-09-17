variable "cluster_name" {
  type        = "string"
  description = "Unique cluster name"
}

# bootkube assets

variable "k8s_domain_name" {
  description = "Master DNS name which resolves to a master instance. Workers and kubeconfig's will communicate with this endpoint (e.g. cluster.example.com)"
  type        = "string"
}

variable "asset_dir" {
  description = "Path to a directory where generated assets should be placed (contains secrets)"
  type        = "string"
}

variable "pod_cidr" {
  description = "CIDR IP range to assign Kubernetes pods"
  type        = "string"
  default     = "10.2.0.0/16"
}

variable "service_cidr" {
  description = <<EOD
CIDR IP range to assign Kubernetes services.
The 1st IP will be reserved for kube_apiserver, the 10th IP will be reserved for kube-dns, the 15th IP will be reserved for self-hosted etcd, and the 200th IP will be reserved for bootstrap self-hosted etcd.
EOD

  type    = "string"
  default = "10.3.0.0/16"
}

variable "master_count" {
  type        = "string"
  default     = "1"
  description = "Number of masters"
}

variable "worker_count" {
  type        = "string"
  default     = "1"
  description = "Number of workers"
}
