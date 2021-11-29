package provider

import (
	"context"
	"github.com/hashicorp/terraform-provider-scaffolding/internal/aidbox"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTokenIntrospector() *schema.Resource {
	return &schema.Resource{
		Description: "TokenIntrospector https://docs.aidbox.app/security-and-access-control-1/auth/access-token-introspection.",
		ReadContext: dataSourceTokenIntrospectorRead,
		Schema:      resourceFullSchema(resourceSchemaTokenIntrospector()),
	}
}

func dataSourceTokenIntrospectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	if id, ok := d.GetOk("name"); ok {
		ti, err := client.GetTokenIntrospector(id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		mapTokenIntrospectorToData(ti, d)
		return nil
	} else {
		return diag.Errorf("No id provided for TokenIntrospector data source")
	}
}
