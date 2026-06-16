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
		CustomizeDiff: customizeAidboxResourceDiff,
		Schema:        resourceFullSchema(resourceSchemaAidboxResource()),
	}
}

func customizeAidboxResourceDiff(_ context.Context, rd *schema.ResourceDiff, _ interface{}) error {
	r1, r2 := rd.GetChange("resource")
	if r1 == "" {
		return nil
	}
	var r1m map[string]interface{}
	var r2m map[string]interface{}
	err := json.Unmarshal([]byte(r1.(string)), &r1m)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(r2.(string)), &r2m)
	if err != nil {
		return err
	}
	id1, _ := r1m["id"]
	id2, _ := r2m["id"]
	if id1 != id2 {
		err = rd.ForceNew("resource")
		if err != nil {
			return err
		}
	}
	rt1, _ := r1m["resourceType"]
	rt2, _ := r2m["resourceType"]
	if rt1 != rt2 {
		err = rd.ForceNew("resource")
		if err != nil {
			return err
		}
	}
	return nil
}

func mapAidboxResourceToData(res *aidbox.GenericResource, data *schema.ResourceData) error {
	data.SetId(res.ResourceTypeAndId)
	// filter the id/meta out here
	var h map[string]any
	err := json.Unmarshal(res.ResourceContent, &h)
	if err != nil {
		return err
	}
	delete(h, "meta")
	if !data.Get("id_assigned").(bool) {
		delete(h, "id")
	}
	resourceContent, err := json.Marshal(h)
	if err != nil {
		return err
	}
	err = data.Set("resource", string(resourceContent))
	if err != nil {
		return err
	}
	return nil
}

func mapAidboxResourceFromData(d *schema.ResourceData) (*aidbox.GenericResource, error) {
	res := &aidbox.GenericResource{}
	res.ResourceTypeAndId = d.Id()
	res.ResourceContent = []byte(d.Get("resource").(string))
	return res, nil
}

func resourceAidboxResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxResourceFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// remember if we assigned the ID ourselves or if we accepted a server-set ID
	// User assigned IDs will continue to be present in the data; server-assigned IDs should not be
	_, err = aidbox.GetResourceTypeAndId(q.ResourceContent)
	if err == nil {
		err = d.Set("id_assigned", true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	res, err := apiClient.CreateGenericResource(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapAidboxResourceToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
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
	err = mapAidboxResourceToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAidboxResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapAidboxResourceFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ti, err := apiClient.UpdateGenericResource(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapAidboxResourceToData(ti, d)
	if err != nil {
		return diag.FromErr(err)
	}
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
	err = mapAidboxResourceToData(res, d)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaAidboxResource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource": {
			Description:      "Aidbox resource content in JSON format",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: jsonDiffSuppressFunc,
		},
		"id_assigned": {
			Description: "Whether an ID was assigned in the original resource or not",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}
