// Create swarm_workers

resource "vultr_instance" "swarm_worker" {
  count             = "${var.worker_instance_count}"
  name              = "${terraform.workspace}-worker-${count.index}"
  region_id         = "${data.vultr_region.my_region.id}"
  plan_id           = "${data.vultr_plan.worker_plan.id}"
  os_id             = "${data.vultr_os.my_os.id}"
  ssh_key_ids       = ["${vultr_ssh_key.nilesh.id}"]
  hostname          = "${terraform.workspace}-worker-${count.index}"
  tag               = "worker"
  private_networking= true
  firewall_group_id = "${vultr_firewall_group.swarm_worker.id}"

  connection {
    type = "ssh"
    user = "root"
  }

  provisioner "remote-exec" {
    inline = [
      "mkdir -p /etc/systemd/system/docker.service.d",
    ]
  }

  provisioner "file" {
    content     = "${data.template_file.docker_conf.rendered}"
    destination = "/etc/systemd/system/docker.service.d/docker.conf"
  }






      provisioner "remote-exec" {
        inline = [
          "echo 'auto ens7' >> /etc/network/interfaces",
          "echo 'iface ens7 inet static' >> /etc/network/interfaces",
          "echo 'address ${self.ipv4_private_address}' >> /etc/network/interfaces",
          "echo 'netmask 255.255.0.0' >> /etc/network/interfaces",
          "echo 'mtu 1450' >> /etc/network/interfaces",
          "ifup ens7"
        ]
      }



  provisioner "file" {
    source      = "scripts/install-docker-ce.sh"
    destination = "/tmp/install-docker-ce.sh"
  }


  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/install-docker-ce.sh",
      "/tmp/install-docker-ce.sh ${var.docker_version}",
      "docker swarm join --token ${data.external.swarm_tokens.result.worker} ${vultr_instance.swarm_manager.0.ipv4_private_address}:2377",
    ]
  }






          #https://www.vultr.com/docs/setup-nfs-share-on-debian
          provisioner "remote-exec" {
            inline = [
              "mkdir /mnt/gluster",
              "apt install -y nfs-common",
              "echo '${vultr_instance.swarm_manager.0.ipv4_private_address}:/mnt/gluster /mnt/gluster  nfs      auto,nofail,noatime,nolock,intr,tcp,actimeo=1800 0 0'   >> /etc/fstab",
              "mount -a",
            ]
          }






  # drain worker on destroy
  provisioner "remote-exec" {
    when = "destroy"
    inline = [
      "docker node update --availability drain ${self.name}",
    ]
    on_failure = "continue"
    connection {
      type = "ssh"
      user = "root"
      host = "${vultr_instance.swarm_manager.0.ipv4_address}"
    }
  }
  # leave swarm on destroy
  provisioner "remote-exec" {
    when = "destroy"
    inline = [
        "docker swarm leave",
    ]
      on_failure = "continue"
  }
  # remove node on destroy
  provisioner "remote-exec" {
    when = "destroy"
      inline = [
        "docker node rm --force ${self.name}",
      ]
      on_failure = "continue"
    connection {
      type = "ssh"
      user = "root"
      host = "${vultr_instance.swarm_manager.0.ipv4_address}"
    }
  }
}
