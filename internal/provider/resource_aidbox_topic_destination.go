package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

const (
	KindTemplate        = "%s-at-least-once"
	KindProfileTemplate = "http://aidbox.app/StructureDefinition/aidboxtopicdestination-%s-at-least-once"
)

// Update is not supported by resource type
func resourceAidboxTopicDestination() *schema.Resource {
	return &schema.Resource{
		Description:   "AidboxTopicDestination https://docs.aidbox.app/modules/topic-based-subscriptions/wip-dynamic-subscriptiontopic-with-destinations",
		CreateContext: resourceAidboxTopicDestinationCreate,
		ReadContext:   resourceAidboxTopicDestinationRead,
		DeleteContext: resourceAidboxTopicDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAidboxTopicDestinationImport,
		},
		Schema: resourceFullSchema(resourceSchemaAidboxTopicDestination()),
	}
}

func resourceSchemaAidboxTopicDestination() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"topic": {
			ForceNew:    true,
			Description: "Reference to the AidboxSubscriptionTopic being subscribed to.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"kind": {
			ForceNew:    true,
			Description: "Defines the destination for sending notifications. Supported values: kafka-at-least-once, kafka-best-effort, webhook-at-least-once, gcp-pubsub-at-least-once, nats-at-least-once, nats-best-effort, amqp-at-least-once, aws-eventbridge, aws-sns, clickhouse, clickhouse-at-least-once, bigquery-at-least-once, data-lakehouse-at-least-once",
			Type:        schema.TypeString,
			Required:    true,
		},
		"content": {
			ForceNew:    true,
			Description: "One of full-resource, id-only or empty",
			Type:        schema.TypeString,
			Required:    true,
		},
		"include_entry_action": {
			ForceNew:    true,
			Description: "When true, each Bundle.entry includes the bundle-entryActionCode extension indicating the CRUD action (create | update | delete) that triggered the notification. Default: true.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"include_version_id": {
			ForceNew:    true,
			Description: "When true, each Bundle.entry includes the bundle-entryVersionId extension containing the resource's meta.versionId at the time of the notification. Default: true.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"parameter": {
			ForceNew:    true,
			Description: "Defines the destination parameters for sending notifications. Parameters are restricted by profiles for each destination.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of a channel-dependent customization parameter",
					},
					"url": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "URL value for the specified parameter name",
					},
					"unsigned_int": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Unsigned integer value for the specified parameter name",
					},
					"string": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "String value for the specified parameter name",
					},
				},
			},
		},
	}
}

func mapAidboxTopicDestinationFromData(data *schema.ResourceData) (*aidbox.AidboxTopicDestination, error) {
	res := &aidbox.AidboxTopicDestination{}
	res.ID = data.Id()
	res.Topic = data.Get("topic").(string)
	res.Content = data.Get("content").(string)
	res.IncludeEntryAction = data.Get("include_entry_action").(bool)
	res.IncludeVersionId = data.Get("include_version_id").(bool)

	// parameter
	rawParameter := data.Get("parameter").([]interface{})
	parameter := make([]aidbox.SubscriptionParameter, len(rawParameter))
	for i, v := range rawParameter {
		raw := v.(map[string]interface{})
		parameter[i] = aidbox.SubscriptionParameter{
			Name:        raw["name"].(string),
			Url:         raw["url"].(string),
			UnsignedInt: raw["unsigned_int"].(int),
			String:      raw["string"].(string),
		}
	}
	res.Parameter = parameter

	// kind
	kind := data.Get("kind").(string)
	res.Kind = fmt.Sprintf(KindTemplate, kind)
	res.Meta = &aidbox.ResourceBaseMeta{
		Profile: []string{fmt.Sprintf(KindProfileTemplate, kind)},
	}

	return res, nil
}

func mapAidboxTopicDestinationToData(res *aidbox.AidboxTopicDestination, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("topic", res.Topic)
	data.Set("content", res.Content)
	data.Set("include_entry_action", res.IncludeEntryAction)
	data.Set("include_version_id", res.IncludeVersionId)

	// parameter
	parameter := make([]interface{}, len(res.Parameter))
	for i, v := range res.Parameter {
		parameter[i] = map[string]interface{}{
			"name":         v.Name,
			"url":          v.Url,
			"unsigned_int": v.UnsignedInt,
			"string":       v.String,
		}
	}
	data.Set("parameter", parameter)
}

func resourceAidboxTopicDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxTopicDestinationFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateAidboxTopicDestination(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxTopicDestinationToData(res, d)
	return nil
}

func resourceAidboxTopicDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAidboxTopicDestination(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapAidboxTopicDestinationToData(res, d)
	return nil
}

func resourceAidboxTopicDestinationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteAidboxTopicDestination(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAidboxTopicDestinationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAidboxTopicDestination(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapAidboxTopicDestinationToData(res, d)
	return []*schema.ResourceData{d}, nil
}
