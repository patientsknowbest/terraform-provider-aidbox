---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "aidbox_user Data Source - terraform-provider-aidbox"
subcategory: ""
description: |-
  User https://docs.aidbox.app/modules/security-and-access-control/readme-1/overview#user
---

# aidbox_user (Data Source)

User https://docs.aidbox.app/modules/security-and-access-control/readme-1/overview#user

## Example Usage

```terraform
data "aidbox_user" "admin_user" {
  id = "admin"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) the ID of the user in aidbox server

### Optional

- `two_factor_enabled` (Boolean) if 2FA is enabled for the user
