// Find the ID for your chosen size/plan.

data "vultr_plan" "manager_plan" {

  filter {
    name   = "price_per_month"
    values = "${var.manager_price_per_month}"
  }

  filter {
    name   = "ram"
    values = "${var.manager_ram}"
  }

}






data "vultr_plan" "worker_plan" {

  filter {
    name   = "price_per_month"
    values = "${var.worker_price_per_month}"
  }

  filter {
    name   = "ram"
    values = "${var.worker_ram}"
  }

}
