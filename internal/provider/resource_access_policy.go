package provider

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "AccessPolicy https://docs.aidbox.app/security-and-access-control-1/security/access-policy.",
		CreateContext: resourceAccessPolicyCreate,
		ReadContext:   resourceAccessPolicyRead,
		UpdateContext: resourceAccessPolicyUpdate,
		DeleteContext: resourceAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAccessPolicyImport,
		},
		Schema: resourceFullSchema(resourceSchemaAccessPolicy()),
	}
}

func mapAccessPolicyToData(res *aidbox.AccessPolicy, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("description", res.Description)
	data.Set("engine", res.Engine.ToString())
	var linkData []interface{}
	for _, ref := range res.Link {
		data := map[string]string{
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
	res, err := apiClient.CreateAccessPolicy(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(res, d)
	return nil
}

func resourceAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAccessPolicy(ctx, d.Id())
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
	ti, err := apiClient.UpdateAccessPolicy(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(ti, d)
	return nil
}

func resourceAccessPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteAccessPolicy(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAccessPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetAccessPolicy(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapAccessPolicyToData(res, d)
	return []*schema.ResourceData{d}, nil
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
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Need to handle empty strings explicitly
				if old == "" && new != "" {
					return false
				}
				if old != "" && new == "" {
					return false
				}
				var oldM, newM map[string]interface{}
				err := json.Unmarshal([]byte(old), &oldM)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal([]byte(new), &newM)
				if err != nil {
					panic(err)
				}
				return reflect.DeepEqual(oldM, newM)
			},
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
