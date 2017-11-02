output "workspace" {
  value = "${terraform.workspace}"
}
output "swarm_manager_token" {
  value = "${data.external.swarm_tokens.result.manager}"
}
output "swarm_worker_token" {
  value = "${data.external.swarm_tokens.result.worker}"
}
output "manager_count" {
  value = ["${vultr_instance.swarm_manager.count}"]
}
output "manager_public_ips" {
  value = ["${vultr_instance.swarm_manager.*.ipv4_address}"]
}
output "manager_private_ips" {
  value = ["${vultr_instance.swarm_manager.*.ipv4_private_address}"]
}
output "worker_count" {
  value = ["${vultr_instance.swarm_worker.count}"]
}
output "worker_public_ips" {
  value = ["${vultr_instance.swarm_worker.*.ipv4_address}"]
}
output "worker_private_ips" {
  value = ["${vultr_instance.swarm_worker.*.ipv4_private_address}"]
}
