package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema:      resourceFullSchema(dataSourceSchemaUser()),
		Description: "User https://docs.aidbox.app/modules/security-and-access-control/readme-1/overview#user",
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetUser(ctx, d.Get("id").(string))
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapUserToData(res, d)
	return nil
}

func mapUserToData(res *aidbox.User, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("two_factor_enabled", res.TwoFactor.Enabled)
}

func dataSourceSchemaUser() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "the ID of the user in aidbox server",
			Type:        schema.TypeString,
			Required:    true,
		},
		"two_factor_enabled": {
			Description: "if 2FA is enabled for the user",
			Type:        schema.TypeBool,
			Optional:    true,
		},
	}
}
