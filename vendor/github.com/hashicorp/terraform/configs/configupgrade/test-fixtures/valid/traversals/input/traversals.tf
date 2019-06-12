locals {
  simple = "${test_instance.foo.bar}"
  splat  = "${test_instance.foo.*.bar}"
  index  = "${test_instance.foo.1.bar}"

  after_simple = "${test_instance.foo.bar.0.baz}"
  after_splat  = "${test_instance.foo.*.bar.0.baz}"
  after_index  = "${test_instance.foo.1.bar.2.baz}"

  non_ident_attr = "${test_instance.foo.bar.1baz}"

  remote_state_output       = "${data.terraform_remote_state.foo.bar}"
  remote_state_attr         = "${data.terraform_remote_state.foo.backend}"
  remote_state_idx_output   = "${data.terraform_remote_state.foo.1.bar}"
  remote_state_idx_attr     = "${data.terraform_remote_state.foo.1.backend}"
  remote_state_splat_output = "${data.terraform_remote_state.foo.*.bar}"
  remote_state_splat_attr   = "${data.terraform_remote_state.foo.*.backend}"

  has_index_should   = "${test_instance.b.0.id}"
  has_index_shouldnt = "${test_instance.c.0.id}"
  no_index_should    = "${test_instance.a.id}"
  no_index_shouldnt  = "${test_instance.c.id}"

  has_index_shouldnt_data = "${data.terraform_remote_state.foo.0.backend}"
}

data "terraform_remote_state" "foo" {
  # This is just here to make sure the schema for this gets loaded to
  # support the remote_state_* checks above.
}

resource "test_instance" "a" {
  count = 1
}

resource "test_instance" "b" {
  count = "${var.count}"
}

resource "test_instance" "c" {
}
