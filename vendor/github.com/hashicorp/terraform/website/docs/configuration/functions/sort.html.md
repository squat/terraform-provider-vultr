---
layout: "functions"
page_title: "sort - Functions - Configuration Language"
sidebar_current: "docs-funcs-collection-sort"
description: |-
  The sort function takes a list of strings and returns a new list with those
  strings sorted lexicographically.
---

# `sort` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`sort` takes a list of strings and returns a new list with those strings
sorted lexicographically.

The sort is in terms of Unicode codepoints, with higher codepoints appearing
after lower ones in the result.

## Examples

```
> sort(["e", "d", "a", "x"])
[
  "a",
  "d",
  "e",
  "x",
]
```
