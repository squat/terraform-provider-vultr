// Create a new SSH key.
resource "vultr_ssh_key" "nilesh" {
  name       = "nilesh"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}
