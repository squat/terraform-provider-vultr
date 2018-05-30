# Vultr Terraform Provider

This is a Terraform provider for Vultr. Find out more about [Vultr](https://www.vultr.com/about/).

[![Build Status](https://travis-ci.org/squat/terraform-provider-vultr.svg?branch=master)](https://travis-ci.org/squat/terraform-provider-vultr)
[![Go Report Card](https://goreportcard.com/badge/github.com/squat/terraform-provider-vultr)](https://goreportcard.com/report/github.com/squat/terraform-provider-vultr)

## Requirements

* A Vultr account and API key
* [Terraform](https://www.terraform.io/downloads.html) 0.9+
* [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

## Usage

Download `terraform-provider-vultr` from the [releases page](https://github.com/squat/terraform-provider-vultr/releases) and follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory,  run `terraform init` to initialize it.

*Note*: in order to build and install the provider from the latest commit on master, run:
```sh
go get -u github.com/squat/terraform-provider-vultr
```

and then register the plugin by symlinking the binary to the [third-party plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins):
```sh
mkdir -p ~/.terraform.d/plugins
ln -s "$GOPATH/bin/terraform-provider-vultr" ~/.terraform.d/plugins/terraform-provider-vultr
```

Set an environment variable containing the Vultr API key:
```
export VULTR_API_KEY=<your-vultr-api-key>
```
*Note*: as an alternative, the API key can be specified in configuration as shown below.

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

## Development

To develop the plugin locally, install the following dependencies:
* [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)
* [Glide](https://github.com/Masterminds/glide#install) (to install and maintain dependencies)
* [glide-vc](https://github.com/sgotti/glide-vc#install) (to clean up dependencies)

To build the plugin run:
```sh
make build
```

To update Go dependencies run:
```sh
make vendor
```
