---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "aidbox_aidbox_topic_destination Resource - terraform-provider-aidbox"
subcategory: ""
description: |-
  AidboxTopicDestination https://docs.aidbox.app/modules/topic-based-subscriptions/wip-dynamic-subscriptiontopic-with-destinations
---

# aidbox_aidbox_topic_destination (Resource)

AidboxTopicDestination https://docs.aidbox.app/modules/topic-based-subscriptions/wip-dynamic-subscriptiontopic-with-destinations



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) One of full-resource, id-only or empty
- `kind` (String) One of kafka, webhook or gcp-pubsub
- `topic` (String) Unique URL of the topic to subscribe on

### Optional

- `parameter` (Block List) Channel-dependent information to send as part of the notification (e.g., HTTP Headers). (see [below for nested schema](#nestedblock--parameter))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--parameter"></a>
### Nested Schema for `parameter`

Required:

- `name` (String) Name of a channel-dependent customization parameter

Optional:

- `string` (String) String value for the specified parameter name
- `unsigned_int` (Number) Unsigned integer value for the specified parameter name
- `url` (String) URL value for the specified parameter name