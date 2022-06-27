package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"
)

func resourceAuthClient() *schema.Resource {
	return &schema.Resource{
		Description:   "Client https://docs.aidbox.app/security-and-access-control-1/auth/basic-auth.",
		CreateContext: resourceAuthClientCreate,
		ReadContext:   resourceAuthClientRead,
		UpdateContext: resourceAuthClientUpdate,
		DeleteContext: resourceAuthClientDelete,
		Schema:        resourceFullSchema(resourceSchemaAuthClient()),
	}
}

func resourceSchemaAuthClient() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"secret": {
			Description: "Client secret used for authentication",
			Type:        schema.TypeString,
			Required:    true,
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

func mapAuthClientToData(res *aidbox.AuthClient, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("secret", res.Secret)
	data.Set("grant_types", res.GrantTypes)
}

func mapAuthClientFromData(d *schema.ResourceData) *aidbox.AuthClient {
	res := &aidbox.AuthClient{}
	res.ID = d.Id()
	res.Secret = d.Get("secret").(string)
	res.GrantTypes = d.Get("grant_types").([]interface{})
	return res
}

func resourceAuthClientCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	q := mapAuthClientFromData(d)
	res, err := client.CreateAuthClient(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapAuthClientToData(res, d)
	return nil
}

func resourceAuthClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	res, err := client.GetAuthClient(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapAuthClientToData(res, d)
	return nil
}

func resourceAuthClientUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	q := mapAuthClientFromData(d)
	ac, err := client.UpdateAuthClient(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapAuthClientToData(ac, d)
	return nil
}

func resourceAuthClientDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	err := client.DeleteAuthClient(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
