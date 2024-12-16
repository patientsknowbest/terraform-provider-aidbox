package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceAidboxSubscriptionTopic() *schema.Resource {
	return &schema.Resource{
		Description:   "AidboxSubscriptionTopic https://docs.aidbox.app/modules/topic-based-subscriptions/wip-dynamic-subscriptiontopic-with-destinations",
		CreateContext: resourceAidboxSubscriptionTopicCreate,
		ReadContext:   resourceAidboxSubscriptionTopicRead,
		UpdateContext: resourceAidboxSubscriptionTopicUpdate,
		DeleteContext: resourceAidboxSubscriptionTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAidboxSubscriptionTopicImport,
		},
		Schema: resourceFullSchema(resourceSchemaAidboxSubscriptionTopic()),
	}
}

func resourceSchemaAidboxSubscriptionTopic() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Description: "Canonical identifier for this search parameter, represented as a URI (globally unique)",
			Type:        schema.TypeString,
			Required:    true,
		},
		"status": {
			Description: "Value of draft | active | retired | unknown, see https://hl7.org/fhir/R4/valueset-publication-status.html",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "active",
		},
		"trigger": {
			Description: "Definition of a trigger for the subscription topic",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "\tKey Data Type, Resource (reference to definition), or relevant definition for this trigger",
					},
				},
			},
		},
	}
}

func mapAidboxSubscriptionTopicFromData(data *schema.ResourceData) (*aidbox.AidboxSubscriptionTopic, error) {
	res := &aidbox.AidboxSubscriptionTopic{}
	res.Url = data.Get("url").(string)
	res.Status = data.Get("status").(string)

	// trigger
	rawTrigger := data.Get("trigger").([]interface{})
	trigger := make([]aidbox.TopicTrigger, len(rawTrigger))
	for i, v := range rawTrigger {
		rawTopicTrigger := v.(map[string]interface{})
		trigger[i] = aidbox.TopicTrigger{
			Resource: rawTopicTrigger["resource"].(string),
		}
	}
	res.Trigger = trigger

	return res, nil
}

func mapAidboxSubscriptionTopicToData(res *aidbox.AidboxSubscriptionTopic, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("url", res.Url)
	data.Set("status", res.Status)

	// trigger
	trigger := make([]interface{}, len(res.Trigger))
	for i, v := range res.Trigger {
		trigger[i] = map[string]interface{}{
			"resource": v.Resource,
		}
	}
	data.Set("trigger", trigger)
}

func resourceAidboxSubscriptionTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxSubscriptionTopicFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateAidboxSubscriptionTopic(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxSubscriptionTopicToData(res, d)
	return nil
}

func resourceAidboxSubscriptionTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAidboxSubscriptionTopic(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapAidboxSubscriptionTopicToData(res, d)
	return nil
}

func resourceAidboxSubscriptionTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxSubscriptionTopicFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateAidboxSubscriptionTopic(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxSubscriptionTopicToData(ac, d)
	return nil
}

func resourceAidboxSubscriptionTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteAidboxSubscriptionTopic(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAidboxSubscriptionTopicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAidboxSubscriptionTopic(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapAidboxSubscriptionTopicToData(res, d)
	return []*schema.ResourceData{d}, nil
}
