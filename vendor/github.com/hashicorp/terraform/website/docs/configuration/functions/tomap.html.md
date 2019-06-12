---
layout: "functions"
page_title: "tomap - Functions - Configuration Language"
sidebar_current: "docs-funcs-conversion-tomap"
description: |-
  The tomap function converts a value to a map.
---

# `tomap` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`tomap` converts its argument to a map value.

Explicit type conversions are rarely necessary in Terraform because it will
convert types automatically where required. Use the explicit type conversion
functions only to normalize types returned in module outputs.

## Examples

```
> tomap({"a" = 1, "b" = 2})
{
  "a" = 1
  "b" = 2
}
```

Since Terraform's concept of a map requires all of the elements to be of the
same type, mixed-typed elements will be converted to the most general type:

```
> tomap({"a" = "foo", "b" = true})
{
  "a" = "foo"
  "b" = "true"
}
```
