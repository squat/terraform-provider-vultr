---
layout: "functions"
page_title: "toset - Functions - Configuration Language"
sidebar_current: "docs-funcs-conversion-toset"
description: |-
  The toset function converts a value to a set.
---

# `toset` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`toset` converts its argument to a set value.

Explicit type conversions are rarely necessary in Terraform because it will
convert types automatically where required. Use the explicit type conversion
functions only to normalize types returned in module outputs.

Pass a _list_ value to `toset` to convert it to a set, which will remove any
duplicate elements and discard the ordering of the elements.

## Examples

```
> toset(["a", "b", "c"])
[
  "a",
  "b",
  "c",
]
```

Since Terraform's concept of a set requires all of the elements to be of the
same type, mixed-typed elements will be converted to the most general type:

```
> toset(["a", "b", 3])
[
  "a",
  "b",
  "3",
]
```

Set collections are unordered and cannot contain duplicate values, so the
ordering of the argument elements is lost and any duplicate values are
coalesced:

```
> toset(["c", "b", "b"])
[
  "b",
  "c",
]
```
