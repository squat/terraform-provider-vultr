// Find the ID for my chosen OS.

data "vultr_os" "my_os" {

  filter {
    name   = "name"
    values = "${var.my_os}"
  }

}
