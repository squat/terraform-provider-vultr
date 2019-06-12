---
layout: "functions"
page_title: "dirname - Functions - Configuration Language"
sidebar_current: "docs-funcs-file-dirname"
description: |-
  The dirname function removes the last portion from a filesystem path.
---

# `dirname` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`dirname` takes a string containing a filesystem path and removes the last
portion from it.

This function works only with the path string and does not access the
filesystem itself. It is therefore unable to take into account filesystem
features such as symlinks.

If the path is empty then the result is `"."`, representing the current
working directory.

The behavior of this function depends on the host platform. On Windows systems,
it uses backslash `\` as the path segment separator. On Unix systems, the slash
`/` is used. The result of this function is normalized, so on a Windows system
any slashes in the given path will be replaced by backslashes before returning.

Referring directly to filesystem paths in resource arguments may cause
spurious diffs if the same configuration is applied from multiple systems or on
different host operating systems. We recommend using filesystem paths only
for transient values, such as the argument to [`file`](./file.html) (where
only the contents are then stored) or in `connection` and `provisioner` blocks.

## Examples

```
> dirname("foo/bar/baz.txt")
foo/bar
```

## Related Functions

* [`basename`](./basename.html) returns _only_ the last portion of a filesystem
  path, discarding the portion that would be returned by `dirname`.
