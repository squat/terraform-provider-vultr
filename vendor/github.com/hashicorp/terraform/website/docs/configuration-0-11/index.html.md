---
layout: "docs"
page_title: "0.11 Configuration Language"
sidebar_current: "docs-conf-old"
description: |-
  Terraform uses text files to describe infrastructure and to set variables. These text files are called Terraform _configurations_ and end in `.tf`. This section talks about the format of these files as well as how they're loaded.
---

# Configuration Language

-> **Note:** This page is about Terraform 0.11 and earlier. For Terraform 0.12
and later, see
[Configuration Language](../configuration/index.html).

Terraform uses text files to describe infrastructure and to set variables.
These text files are called Terraform _configurations_ and end in
`.tf`. This section talks about the format of these files as well as
how they're loaded.

The format of the configuration files are able to be in two formats:
Terraform format and JSON. The Terraform format is more human-readable,
supports comments, and is the generally recommended format for most
Terraform files. The JSON format is meant for machines to create,
modify, and update, but can also be done by Terraform operators if
you prefer. Terraform format ends in `.tf` and JSON format ends in
`.tf.json`.

Click a sub-section in the navigation to the left to learn more about
Terraform configuration.
