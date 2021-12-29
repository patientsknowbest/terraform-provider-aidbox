package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-scaffolding/internal/aidbox"
)

func resourceSchemaTokenIntrospector() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

//const (
//	aidboxTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
//)

func mapResourceBaseToData(v *aidbox.ResourceBase, data *schema.ResourceData) {
	if v.ID != "" {
		data.SetId(v.ID)
	}
	//if v.Meta != nil {
	//	meta := map[string]interface{}{}
	//	if v.Meta.VersionId != "" {
	//		meta["version_id"] = v.Meta.VersionId
	//	}
	//	if v.Meta.CreatedAt != nil {
	//		meta["created_at"] = v.Meta.CreatedAt.Format(aidboxTimeFormat)
	//	}
	//	if v.Meta.LastUpdated != nil {
	//		meta["last_updated"] = v.Meta.LastUpdated.Format(aidboxTimeFormat)
	//	}
	//	data.Set("meta", meta)
	//}
}

func mapResourceBaseFromData(d *schema.ResourceData) aidbox.ResourceBase {
	res := aidbox.ResourceBase{}
	res.ID = d.Id()
	//if v, ok := d.GetOk("meta"); ok {
	//	meta := v.([]interface{})[0].(map[string]interface{}) // Ugly
	//	mm := &aidbox.ResourceBaseMeta{}
	//	if vv, ok := meta["created_at"]; ok {
	//		if vvv, err := time.Parse(aidboxTimeFormat, vv.(string)); err != nil {
	//			mm.CreatedAt = &vvv
	//		}
	//	}
	//	if vv, ok := meta["last_updated"]; ok {
	//		if vvv, err := time.Parse(aidboxTimeFormat, vv.(string)); err != nil {
	//			mm.LastUpdated = &vvv
	//		}
	//	}
	//	if vv, ok := meta["version_id"]; ok {
	//		mm.VersionId = vv.(string)
	//	}
	//	res.Meta = mm
	//}
	return res
}

func mapTokenIntrospectorToData(v *aidbox.TokenIntrospector, data *schema.ResourceData) {
	mapResourceBaseToData(&v.ResourceBase, data)
	if v.TokenIntrospectionEndpoint != nil {
		data.Set("introspection_endpoint", map[string]interface{}{
			"authorization": v.TokenIntrospectionEndpoint.Authorization,
			"url":           v.TokenIntrospectionEndpoint.URL,
		})
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
	client := meta.(*aidbox.Client)
	q := mapTokenIntrospectorFromData(d)
	res, err := client.CreateTokenIntrospector(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(res, d)
	return nil
}

func resourceTokenIntrospectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	res, err := client.GetTokenIntrospector(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(res, d)
	return nil
}

func resourceTokenIntrospectorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*aidbox.Client)
	q := mapTokenIntrospectorFromData(d)
	ti, err := client.UpdateTokenIntrospector(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(ti, d)
	return nil
}

func resourceTokenIntrospectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*aidbox.Client)
	err := client.DeleteTokenIntrospector(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
