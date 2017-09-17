# Secure copy bootkube assets to ONE master and start bootkube to perform
# one-time self-hosted cluster bootstrapping.
resource "null_resource" "bootkube-start" {
  depends_on = ["module.bootkube", "vultr_instance.masters"]

  connection {
    type    = "ssh"
    host    = "${vultr_instance.masters.0.ipv4_address}"
    user    = "core"
    timeout = "15m"
  }

  provisioner "file" {
    source      = "${var.asset_dir}"
    destination = "$HOME/assets"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv /home/core/assets /opt/bootkube",
      "sudo systemctl start bootkube",
    ]
  }
}
