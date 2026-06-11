package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceAidboxResource() *schema.Resource {
	return &schema.Resource{
		Description: "AidboxResource generic resource for any aidbox resource. " +
			"Escape hatch for controlling resources not natively supported by this plugin",
		CreateContext: resourceAidboxResourceCreate,
		ReadContext:   resourceAidboxResourceRead,
		UpdateContext: resourceAidboxResourceUpdate,
		DeleteContext: resourceAidboxResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAidboxResourceImport,
		},
		Schema: resourceFullSchema(resourceSchemaAidboxResource()),
	}
}

func mapAidboxResourceToData(res *aidbox.GenericResource, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("resource", string(res.ResourceContent))
}

func mapAidboxResourceFromData(d *schema.ResourceData) *aidbox.GenericResource {
	res := &aidbox.GenericResource{}
	res.ID = d.Id()
	res.ResourceContent = json.RawMessage(d.Get("resource").(string))
	return res
}

func resourceAidboxResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapAidboxResourceFromData(d)
	res, err := apiClient.CreateGenericResource(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxResourceToData(res, d)
	return nil
}

func resourceAidboxResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetGenericResource(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapAidboxResourceToData(res, d)
	return nil
}

func resourceAidboxResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapAidboxResourceFromData(d)
	ti, err := apiClient.UpdateGenericResource(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapAidboxResourceToData(ti, d)
	return nil
}

func resourceAidboxResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteGenericResource(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAidboxResourceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetGenericResource(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapAidboxResourceToData(res, d)
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaAidboxResource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource": {
			Description:      "Aidbox resource content in JSON format",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: jsonDiffSuppressMetaFunc,
		},
	}
}
