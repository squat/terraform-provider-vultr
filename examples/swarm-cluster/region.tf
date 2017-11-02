// Find the ID of the chosen region.

data "vultr_region" "my_region" {
  filter {
    name   = "name"
    values = "${var.region}"
  }
}
