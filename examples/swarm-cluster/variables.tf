variable "region" {
  default = ["Seattle"]
}

variable "my_os" {
  default = ["Debian 9 x64 (stretch)"]
}

variable "manager_price_per_month" {
  default =  ["10.00"]
}
variable "manager_ram" {
  default =  ["2048"]
}

variable "worker_price_per_month" {
  default =  ["5.00"]
}
variable "worker_ram" {
  default = ["1024"]
}

variable "manager_instance_count" {
  default = 1
}

variable "worker_instance_count" {
  default = 1
}


//apt-cache madison docker-ce
variable "docker_version" {
  default = "17.09.0~ce-0~debian"
}
variable "docker_api_ip" {
  default = "127.0.0.1"
}

data "template_file" "docker_conf" {
  template = "${file("conf/docker.tpl")}"

  vars {
    ip = "${var.docker_api_ip}"
  }
}
