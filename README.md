Vultr Terraform Provider
========================

This is a [Terraform][tf] provider for the Vultr cloud. Find out more about
[Vultr][vultr].

[![Build Status][build-status-img]][build-status]
[![Go Report Card][go-report-card-img]][go-report-card]

Requirements
------------

- A Vultr account and API key
- [Terraform](https://www.terraform.io/downloads.html) 0.10.x
- [Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

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
your machine (version 1.9+ is *required*). You'll also need to correctly setup a
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

[build-status]: https://travis-ci.org/squat/terraform-provider-vultr
[build-status-img]: https://travis-ci.org/squat/terraform-provider-vultr.svg?branch=master
[go-report-card]: https://goreportcard.com/report/github.com/squat/terraform-provider-vultr
[go-report-card-img]: https://goreportcard.com/badge/github.com/squat/terraform-provider-vultr
[golang]: https://www.golang.org/
[tf]: https://www.terraform.io/
[tf-plugin]: https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin
[vultr]: https://www.vultr.com/about/
