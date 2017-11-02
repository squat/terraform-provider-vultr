// Create swarm_manager

resource "vultr_instance" "swarm_manager" {
  count             = "${var.manager_instance_count}"
  name              = "${terraform.workspace}-manager-${count.index}"
  region_id         = "${data.vultr_region.my_region.id}"
  plan_id           = "${data.vultr_plan.manager_plan.id}"
  os_id             = "${data.vultr_os.my_os.id}"
  ssh_key_ids       = ["${data.vultr_ssh_key.nilesh.id}"]
  hostname          = "${terraform.workspace}-manager-${count.index}"
  tag               = "manager"
  private_networking= true
  firewall_group_id = "${vultr_firewall_group.swarm_manager.id}"


    connection {
      type = "ssh"
      user = "root"
    }

    provisioner "remote-exec" {
      inline = [
        "mkdir -p /etc/systemd/system/docker.service.d",
      ]
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
        "docker swarm init --advertise-addr ${self.ipv4_private_address}",
      ]
    }


            #https://www.vultr.com/docs/setup-nfs-share-on-debian
            provisioner "remote-exec" {
              inline = [
                "apt-get install -y nfs-kernel-server nfs-common",
                "mkdir /mnt/gluster",
                "chown nobody:nogroup /mnt/gluster",
                "chmod 755 /mnt/gluster",
                "echo '/mnt/gluster   ${self.ipv4_private_address}/24(rw,no_root_squash,sync,no_subtree_check)' >> /etc/exports",
                "chmod -R o+w /mnt/gluster/",
                "service nfs-kernel-server restart",
              ]
            }





}
