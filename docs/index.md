---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kafka Provider"
subcategory: ""
description: |-
---

Please, please, please, for the love of everything you hold dear, **DO NOT** use this in production!!!

# kafka Provider

I mean, seriously.  I developed this against a dev docker container.  It did what I needed it to do for *testing*.

## Example Usage

```hcl
provider "kafka" {
    bootstrap_servers = ["localhost:9092"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bootstrap_servers` (List of String) a list of kafka brokers