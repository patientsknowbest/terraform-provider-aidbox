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

func resourceAidboxTopicDestination() *schema.Resource {
	return &schema.Resource{
		Description:   "AidboxTopicDestination https://docs.aidbox.app/modules/topic-based-subscriptions/wip-dynamic-subscriptiontopic-with-destinations",
		CreateContext: resourceAidboxTopicDestinationCreate,
		ReadContext:   resourceAidboxTopicDestinationRead,
		UpdateContext: resourceAidboxTopicDestinationUpdate,
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
			Description: "Unique URL of the topic to subscribe on",
			Type:        schema.TypeString,
			Required:    true,
		},
		"kind": {
			Description: "One of kafka, webhook or gcp-pubsub",
			Type:        schema.TypeString,
			Required:    true,
		},
		"content": {
			Description: "One of full-resource, id-only or empty",
			Type:        schema.TypeString,
			Required:    true,
		},
		"parameter": {
			Description: "Channel-dependent information to send as part of the notification (e.g., HTTP Headers).",
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
	res.Topic = data.Get("topic").(string)
	res.Content = data.Get("content").(string)

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

func resourceAidboxTopicDestinationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxTopicDestinationFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateAidboxTopicDestination(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxTopicDestinationToData(ac, d)
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
