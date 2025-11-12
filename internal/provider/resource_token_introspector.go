package provider

import (
	"context"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceTokenIntrospector() *schema.Resource {
	return &schema.Resource{
		Description:   "TokenIntrospector https://docs.aidbox.app/security-and-access-control-1/auth/access-token-introspection.",
		CreateContext: resourceTokenIntrospectorCreate,
		ReadContext:   resourceTokenIntrospectorRead,
		UpdateContext: resourceTokenIntrospectorUpdate,
		DeleteContext: resourceTokenIntrospectorDelete,
		Schema:        resourceFullSchema(resourceSchemaTokenIntrospector()),
	}
}

func resourceSchemaTokenIntrospector() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Description: "Type of token introspector. One of (opaque|jwt)",
			Type:        schema.TypeString,
			Required:    true,
			ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
				_, err := aidbox.ParseTokenIntrospectorType(i.(string))
				if err != nil {
					return diag.FromErr(err)
				}
				return nil
			},
		},
		"introspection_endpoint": {
			// This description is used by the documentation generator and the language server.
			Description: "Configuration for introspecting opaque access tokens.",
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"authorization": {
						Description: "Authorization header value.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"url": {
						Description: "URL of the introspection endpoint.",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
		},
		"jwks_uri": {
			Description: "Location of JWKS public key information for validating JWT tokens",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"jwt": {
			Description: "Configuration for validating jwt type access tokens",
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"iss": {
						Description: "The issuer of the JWT",
						Type:        schema.TypeString,
						Required:    true,
					},
					"secret": {
						Description: "The secret used to sign the JWT",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}
func mapTokenIntrospectorToData(v *aidbox.TokenIntrospector, data *schema.ResourceData) {
	mapResourceBaseToData(&v.ResourceBase, data)
	data.Set("type", v.Type.ToString())

	if v.TokenIntrospectionEndpoint != nil {
		ti := map[string]interface{}{
			"authorization": v.TokenIntrospectionEndpoint.Authorization,
			"url":           v.TokenIntrospectionEndpoint.URL,
		}
		data.Set("introspection_endpoint", []interface{}{ti})
	}
	if v.JWKSURI != "" {
		data.Set("jwks_uri", v.JWKSURI)
	}
	if v.TokenIntrospectorJWT != nil {
		jwt := map[string]interface{}{
			"iss":    v.TokenIntrospectorJWT.ISS,
			"secret": v.TokenIntrospectorJWT.Secret,
		}
		data.Set("jwt", []interface{}{jwt})
	}
}

func mapTokenIntrospectorFromData(d *schema.ResourceData) *aidbox.TokenIntrospector {
	vv := &aidbox.TokenIntrospector{
		ResourceBase: mapResourceBaseFromData(d),
	}

	if tt, err := aidbox.ParseTokenIntrospectorType(d.Get("type").(string)); err == nil {
		vv.Type = tt
	} else {
		log.Panicf("Error parsing token introspector type %v", err)
	}

	if v, ok := d.GetOk("introspection_endpoint"); ok {
		introspectionEndpointData := v.([]interface{})[0].(map[string]interface{}) // Ugly
		vv.TokenIntrospectionEndpoint = &aidbox.TokenIntrospectionEndpoint{
			Authorization: introspectionEndpointData["authorization"].(string),
			URL:           introspectionEndpointData["url"].(string),
		}
	}
	if v, ok := d.GetOk("jwks_uri"); ok {
		vv.JWKSURI = v.(string)
	}
	if v, ok := d.GetOk("jwt"); ok {
		jwtData := v.([]interface{})[0].(map[string]interface{}) // Ugly
		vv.TokenIntrospectorJWT = &aidbox.TokenIntrospectorJWT{
			ISS:    jwtData["iss"].(string),
			Secret: jwtData["secret"].(string),
		}
	}

	return vv
}

func resourceTokenIntrospectorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapTokenIntrospectorFromData(d)
	res, err := apiClient.CreateTokenIntrospector(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(res, d)
	return nil
}

func resourceTokenIntrospectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetTokenIntrospector(ctx, d.Id())

	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(res, d)
	return nil
}

func resourceTokenIntrospectorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapTokenIntrospectorFromData(d)
	ti, err := apiClient.UpdateTokenIntrospector(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(ti, d)
	return nil
}

func resourceTokenIntrospectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteTokenIntrospector(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
