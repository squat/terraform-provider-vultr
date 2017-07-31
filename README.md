# Vultr Terraform Provider

This is a Terraform provider for the Vultr cloud. Find out more about [Vultr](https://www.vultr.com/about/).

[![Build Status](https://travis-ci.org/squat/terraform-provider-vultr.svg?branch=master)](https://travis-ci.org/squat/terraform-provider-vultr)
[![Go Report Card](https://goreportcard.com/badge/github.com/squat/terraform-provider-vultr)](https://goreportcard.com/report/github.com/squat/terraform-provider-vultr)

## Requirements

* A Vultr account and API key
* [Terraform](https://www.terraform.io/downloads.html) 0.9+
* [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

## Usage

Download `terraform-provider-vultr` and install the plugin binary on the filesystem:
```sh
go get -u github.com/squat/terraform-provider-vultr
```

Register the plugin in `~/.terraformrc`:
```sh
providers {
  vultr = "/path/to/terraform-provider-vultr"
}
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

// Create a Vultr virtual machine.
resource "vultr_instance" "example" {
  name              = "example"
  region_id         = 12                                   // Silicon Valley
  plan_id           = 201                                  // $5
  os_id             = 179                                  // CoreOS Container Linux stable
  ssh_keys          = ["squat"]
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
* [Dep](https://github.com/golang/dep#setup) (to install and maintain dependencies)

To build the plugin run:
```sh
make build
```

To update Go dependencies run:
```sh
make vendor
```
