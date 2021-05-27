---
layout: "sap"
page_title: "Provider: SAP"
description: |-
  The SAP provider is used to interact with the many resources supported by SAP. The provider needs to be configured with the proper credentials before it can be used.
---

# SAP Provider

The SAP provider is used to interact with the
many resources supported by SAP. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    sap = {
      source  = "nnicora/sap"
      version = "~> 0.0.38"
    }
  }
}
```