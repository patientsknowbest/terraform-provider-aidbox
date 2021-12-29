package provider

import (
	"context"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"

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
	ti, err := client.GetTokenIntrospector(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	mapTokenIntrospectorToData(ti, d)
	return nil
}
