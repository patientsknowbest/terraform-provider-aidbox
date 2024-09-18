package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceClient() *schema.Resource {
	return &schema.Resource{
		Description:   "Client https://docs.aidbox.app/security-and-access-control-1/auth/basic-auth.",
		CreateContext: resourceClientCreate,
		ReadContext:   resourceClientRead,
		UpdateContext: resourceClientUpdate,
		DeleteContext: resourceClientDelete,
		Importer:      &schema.ResourceImporter{
			StateContext: resourceClientImport,
		},
		Schema:        resourceFullSchema(resourceSchemaClient()),
	}
}

func resourceSchemaClient() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Client ID used for authentication",
			Type:        schema.TypeString,
			Required:    true,
		},
		"secret": {
			Description: "Client secret used for authentication",
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
		},
		"grant_types": {
			Description: "Grant type used for authentication (basic)",
			Type:        schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
			MinItems: 1,
		},
	}
}

func mapClientToData(res *aidbox.Client, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("name", res.ID)
	data.Set("secret", res.Secret)
	var types []interface{}
	for _, gt := range res.GrantTypes {
		types = append(types, gt.ToString())
	}
	data.Set("grant_types", types)
}

func mapClientFromData(d *schema.ResourceData) (*aidbox.Client, error) {
	res := &aidbox.Client{}
	res.ID = d.Get("name").(string)
	res.Secret = d.Get("secret").(string)
	types := d.Get("grant_types").([]interface{})
	grantTypes := []aidbox.GrantType{}
	for _, t := range types {
		gt, err := aidbox.ParseGrantType(fmt.Sprint(t))
		if err != nil {
			return nil, err
		}
		grantTypes = append(grantTypes, gt)
	}
	res.GrantTypes = grantTypes
	return res, nil
}

func resourceClientCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapClientFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateClient(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapClientToData(res, d)
	return nil
}

func resourceClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetClient(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapClientToData(res, d)
	return nil
}

func resourceClientUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapClientFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateClient(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapClientToData(ac, d)
	return nil
}

func resourceClientDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteClient(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceClientImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetClient(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return nil, err
	}
	mapClientToData(res, d)
	return []*schema.ResourceData{d}, nil
}