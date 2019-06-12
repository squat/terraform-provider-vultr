---
layout: "functions"
page_title: "values - Functions - Configuration Language"
sidebar_current: "docs-funcs-collection-values"
description: |-
  The values function returns a list of the element values in a given map.
---

# `values` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`values` takes a map and returns a list containing the values of the elements
in that map.

The values are returned in lexicographical order by their corresponding _keys_,
so the values will be returned in the same order as their keys would be
returned from [`keys`](./keys.html).

## Examples

```
> values({a=3, c=2, d=1})
[
  3,
  2,
  1,
]
```

## Related Functions

* [`keys`](./keys.html) returns a list of the _keys_ from a map.
