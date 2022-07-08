package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"
	"log"
)

func resourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "AccessPolicy https://docs.aidbox.app/security-and-access-control-1/security/access-policy.",
		CreateContext: resourceAccessPolicyCreate,
		ReadContext:   resourceAccessPolicyRead,
		UpdateContext: resourceAccessPolicyUpdate,
		DeleteContext: resourceAccessPolicyDelete,
		Schema:        resourceFullSchema(resourceSchemaAccessPolicy()),
	}
}

func mapAccessPolicyToData(res *aidbox.AccessPolicy, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("description", res.Description)
	data.Set("engine", res.Engine.ToString())
	var linkData []interface{}
	for _, ref := range res.Link {
		data := map[string]interface{}{
			"resource_id":   ref.ResourceId,
			"resource_type": ref.ResourceType,
		}
		linkData = append(linkData, data)
	}
	data.Set("link", linkData)
	if string(res.Schema) != "" {
		data.Set("schema", string(res.Schema))
	}
}

func mapAccessPolicyFromData(d *schema.ResourceData) *aidbox.AccessPolicy {
	res := &aidbox.AccessPolicy{}
	res.ID = d.Id()
	res.Description = d.Get("description").(string)
	e, err := aidbox.ParseAccessPolicyEngine(d.Get("engine").(string))
	if err != nil {
		log.Panicln(err)
	}
	res.Engine = e
	if v, ok := d.GetOk("link"); ok {
		references := []aidbox.Reference{}
		for _, data := range v.([]interface{}) {
			linkData := data.(map[string]interface{})
			ref := aidbox.Reference{
				ResourceId:   linkData["resource_id"].(string),
				ResourceType: linkData["resource_type"].(string),
			}
			references = append(references, ref)
		}
		res.Link = references
	}
	if vv, ok := d.GetOk("schema"); ok {
		res.Schema = json.RawMessage(vv.(string))
	}
	return res
}

func resourceAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapAccessPolicyFromData(d)
	res, err := apiClient.CreateAccessPolicy(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(res, d)
	return nil
}

func resourceAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAccessPolicy(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(res, d)
	return nil
}

func resourceAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapAccessPolicyFromData(d)
	ti, err := apiClient.UpdateAccessPolicy(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(ti, d)
	return nil
}

func resourceAccessPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteAccessPolicy(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSchemaAccessPolicy() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "Description of access policy for human users.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"engine": {
			Description: "The engine which is used to evaluate this policy. One of (json-schema|allow)",
			Type:        schema.TypeString,
			Required:    true,
		},
		"schema": {
			Description: "JSON-schema policy to be evaluated. Used only if engine is json-schema",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"link": {
			Description: "The actor to allow access. Used only if engine is allow.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource_id": {
						Description: "The ID of the referenced resource",
						Type:        schema.TypeString,
						Required:    true,
					},
					"resource_type": {
						Description: "The type of the referenced resource (Client)",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
			},
		},
	}
}
