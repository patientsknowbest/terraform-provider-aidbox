package provider

import (
	"context"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "AccessPolicy https://docs.aidbox.app/security-and-access-control-1/security/access-policy.",
		ReadContext: dataSourceAccessPolicyRead,
		Schema:      resourceFullSchema(resourceSchemaAccessPolicy()),
	}
}

func dataSourceAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*aidbox.Client)
	ti, err := client.GetAccessPolicy(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	mapAccessPolicyToData(ti, d)
	return nil
}
