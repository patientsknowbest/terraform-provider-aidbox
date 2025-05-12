package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
	"log"
)

func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Description:   "IdentityProvider https://docs.aidbox.app/modules/security-and-access-control/set-up-external-identity-provider.",
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityProviderImport,
		},
		Schema: resourceFullSchema(resourceSchemaIdentityProvider()),
	}
}

func resourceSchemaIdentityProvider() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Description: "Title of the identity provider.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"client": {
			Description: "Authentication of the OAuth Provider.",
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "id of the client you registered in OAuth Provider API.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"secret": {
						Description: "secret of the client you registered in OAuth Provider API.",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
		},
		"system": {
			Description: "Adds identifier for the created user with this system.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"authorize_endpoint": {
			Description: "OAuth Provider authorization endpoint.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"token_endpoint": {
			Description: "OAuth Provider access token endpoint.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"userinfo_source": {
			Description: "One of (id-token|userinfo-endpoint). If `id-token`, then `user.data` is populated with the `id_token.claims` value. Otherwise request to the `userinfo_endpoint` is performed to get user details.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"userinfo_endpoint": {
			Description: "OAuth Provider user profile endpoint.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"scopes": {
			Description: "Array of scopes for which you request access from user.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
func mapIdentityProviderToData(v *aidbox.IdentityProvider, data *schema.ResourceData) {
	mapResourceBaseToData(&v.ResourceBase, data)

	data.Set("title", v.Title)
	data.Set("system", v.System)
	data.Set("authorize_endpoint", v.AuthorizeEndpoint)
	data.Set("token_endpoint", v.TokenEndpoint)
	data.Set("userinfo_source", v.UserinfoSource.ToString())
	data.Set("userinfo_endpoint", v.UserinfoEndpoint)
	data.Set("scopes", v.Scopes)

	if v.Client != nil {
		client := map[string]interface{}{
			"id":     v.Client.ID,
			"secret": v.Client.Secret,
		}
		data.Set("client", []interface{}{client})
	}
}

func mapIdentityProviderFromData(d *schema.ResourceData) *aidbox.IdentityProvider {
	vv := &aidbox.IdentityProvider{
		ResourceBase: mapResourceBaseFromData(d),
	}

	vv.Title = d.Get("title").(string)
	vv.System = d.Get("system").(string)
	vv.AuthorizeEndpoint = d.Get("authorize_endpoint").(string)
	vv.TokenEndpoint = d.Get("token_endpoint").(string)
	vv.UserinfoEndpoint = d.Get("userinfo_endpoint").(string)

	scopes := d.Get("scopes").([]interface{})
	for _, scope := range scopes {
		vv.Scopes = append(vv.Scopes, scope.(string))
	}

	if tt, err := aidbox.ParseUserinfoSource(d.Get("userinfo_source").(string)); err == nil {
		vv.UserinfoSource = tt
	} else {
		log.Panicln(err)
	}

	if v, ok := d.GetOk("client"); ok {
		clientData := v.([]interface{})[0].(map[string]interface{}) // Ugly
		vv.Client = &aidbox.IdentityProviderClient{
			ID:     clientData["id"].(string),
			Secret: clientData["secret"].(string),
		}
	}

	return vv
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapIdentityProviderFromData(d)
	res, err := apiClient.CreateIdentityProvider(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapIdentityProviderToData(res, d)
	return nil
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetIdentityProvider(ctx, d.Id())

	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapIdentityProviderToData(res, d)
	return nil
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapIdentityProviderFromData(d)
	ti, err := apiClient.UpdateIdentityProvider(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapIdentityProviderToData(ti, d)
	return nil
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteIdentityProvider(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIdentityProviderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetIdentityProvider(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapIdentityProviderToData(res, d)
	return []*schema.ResourceData{d}, nil
}
