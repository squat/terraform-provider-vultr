Vultr Terraform Provider
========================

This is a [Terraform][tf] provider for the Vultr cloud. Find out more about
[Vultr][vultr].

[![Build Status][build-status-img]][build-status]
[![Go Report Card][go-report-card-img]][go-report-card]

Requirements
------------

- A Vultr account and API key
- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/squat/terraform-provider-vultr`

```sh
$ mkdir -p $GOPATH/src/github.com/squat; cd $GOPATH/src/github.com/squat
$ git clone git@github.com:squat/terraform-provider-vultr
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/squat/terraform-provider-vultr
$ make build
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it][tf-plugin]
as a plugin. After placing it into your plugins directory,  run `terraform init`
to initialize it.

The Vultr API key can be provided as environment variable

```bash
export VULTR_API_KEY=<your-vultr-api-key>
```

or on the Terraform provider configuration:

```hcl
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go][golang] installed on
your machine (version 1.10+ is *required*). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin`
to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-vultr
...
```

In order to test the provider, you can simply run `make test`.

*Note:* Make sure no `VULTR_API_KEY` variables is set.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Examples

```tf
// Configure the Vultr provider. 
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the ID of the Silicon Valley region.
data "vultr_region" "silicon_valley" {
  filter {
    name   = "name"
    values = ["Silicon Valley"]
  }
}

// Find the ID for CoreOS Container Linux.
data "vultr_os" "container_linux" {
  filter {
    name   = "family"
    values = ["coreos"]
  }
}

// Find the ID for a starter plan.
data "vultr_plan" "starter" {
  filter {
    name   = "price_per_month"
    values = ["5.00"]
  }

  filter {
    name   = "ram"
    values = ["1024"]
  }
}

// Find the ID of an existing SSH key.
data "vultr_ssh_key" "squat" {
  filter {
    name   = "name"
    values = ["squat"]
  }
}

// Create a Vultr virtual machine.
resource "vultr_instance" "example" {
  name              = "example"
  region_id         = "${data.vultr_region.silicon_valley.id}"
  plan_id           = "${data.vultr_plan.starter.id}"
  os_id             = "${data.vultr_os.container_linux.id}"
  ssh_key_ids       = ["${data.vultr_ssh_key.squat.id}"]
  hostname          = "example"
  tag               = "container-linux"
  firewall_group_id = "${vultr_firewall_group.example.id}"
}

// Create a new firewall group.
resource "vultr_firewall_group" "example" {
  description = "example group"
}

// Add a firewall rule to the group allowing SSH access.
resource "vultr_firewall_rule" "ssh" {
  firewall_group_id = "${vultr_firewall_group.example.id}"
  cidr_block        = "0.0.0.0/0"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
}
```


[build-status]: https://travis-ci.org/squat/terraform-provider-vultr
[build-status-img]: https://travis-ci.org/squat/terraform-provider-vultr.svg?branch=master
[go-report-card]: https://goreportcard.com/report/github.com/squat/terraform-provider-vultr
[go-report-card-img]: https://goreportcard.com/badge/github.com/squat/terraform-provider-vultr
[golang]: https://www.golang.org/
[tf]: https://www.terraform.io/
[tf-plugin]: https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin
[vultr]: https://www.vultr.com/about/