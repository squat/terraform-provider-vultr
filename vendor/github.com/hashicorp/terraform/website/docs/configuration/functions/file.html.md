---
layout: "functions"
page_title: "file - Functions - Configuration Language"
sidebar_current: "docs-funcs-file-file-x"
description: |-
  The file function reads the contents of the file at the given path and
  returns them as a string.
---

# `file` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`file` reads the contents of a file at the given path and returns them as
a string.

```hcl
file(path)
```

Strings in the Terraform language are sequences of Unicode characters, so
this function will interpret the file contents as UTF-8 encoded text and
return the resulting Unicode characters. If the file contains invalid UTF-8
sequences then this function will produce an error.

This function can be used only with files that already exist on disk
at the beginning of a Terraform run. Functions do not participate in the
dependency graph, so this function cannot be used with files that are generated
dynamically during a Terraform operation. We do not recommend using dynamic
local files in Terraform configurations, but in rare situations where this is
necessary you can use
[the `local_file` data source](/docs/providers/local/d/file.html)
to read files while respecting resource dependencies.

## Examples

```
> file("${path.module}/hello.txt")
Hello World
```

## Related Functions

* [`filebase64`](./filebase64.html) also reads the contents of a given file,
  but returns the raw bytes in that file Base64-encoded, rather than
  interpreting the contents as UTF-8 text.
* [`fileexists`](./fileexists.html) determines whether a file exists
  at a given path.
* [`templatefile`](./templatefile.html) renders using a file from disk as a
  template.
