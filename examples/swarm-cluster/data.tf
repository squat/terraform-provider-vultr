data "external" "swarm_tokens" {
  program = ["./scripts/get-join-tokens.sh"]
  query = {
    host = "${vultr_instance.swarm_manager.0.ipv4_address}"
  }

    depends_on = ["vultr_instance.swarm_manager"]
}
