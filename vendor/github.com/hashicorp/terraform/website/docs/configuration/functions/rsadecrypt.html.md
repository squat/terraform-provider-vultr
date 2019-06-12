---
layout: "functions"
page_title: "rsadecrypt - Functions - Configuration Language"
sidebar_current: "docs-funcs-crypto-rsadecrypt"
description: |-
  The rsadecrypt function decrypts an RSA-encrypted message.
---

# `rsadecrypt` Function

-> **Note:** This page is about Terraform 0.12 and later. For Terraform 0.11 and
earlier, see
[0.11 Configuration Language: Interpolation Syntax](../../configuration-0-11/interpolation.html).

`rsadecrypt` decrypts an RSA-encrypted ciphertext, returning the corresponding
cleartext.

```hcl
rsadecrypt(ciphertext, privatekey)
```

`ciphertext` must be a base64-encoded representation of the ciphertext, using
the PKCS #1 v1.5 padding scheme. Terraform uses the "standard" Base64 alphabet
as defined in [RFC 4648 section 4](https://tools.ietf.org/html/rfc4648#section-4).

`privatekey` must be a PEM-encoded RSA private key that is not itself
encrypted.

Terraform has no corresponding function for _encrypting_ a message. Use this
function to decrypt ciphertexts returned by remote services using a keypair
negotiated out-of-band.

## Examples

```
> rsadecrypt(filebase64("${path.module}/ciphertext"), file("privatekey.pem"))
Hello, world!
```
