package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceGcpServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Aidbox GcpServiceAccount is a proprietary custom resource used to store Google Cloud Platform service account credentials for workload identity.",
		CreateContext: resourceGcpServiceAccountCreateOrUpdate,
		ReadContext:   resourceGcpServiceAccountRead,
		UpdateContext: resourceGcpServiceAccountCreateOrUpdate,
		DeleteContext: resourceGcpServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGcpServiceAccountImport,
		},
		Schema: resourceFullSchema(resourceSchemaGcpServiceAccount()),
	}
}

func resourceSchemaGcpServiceAccount() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Computer friendly name of the GCP Service Account.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"service_account_email": {
			Description: "The email address of the GCP service account.",
			Type:        schema.TypeString,
			Required:    true,
		},
	}
}

func mapGcpServiceAccountFromData(data *schema.ResourceData) (*aidbox.GcpServiceAccount, error) {
	res := &aidbox.GcpServiceAccount{}
	res.ResourceType = "GcpServiceAccount"
	res.ID = data.Get("name").(string)

	if v, ok := data.GetOk("service_account_email"); ok {
		res.ServiceAccountEmail = v.(string)
	}

	return res, nil
}

func mapGcpServiceAccountToData(res *aidbox.GcpServiceAccount, data *schema.ResourceData) error {
	data.SetId(res.ID)
	data.Set("name", res.ID)
	data.Set("service_account_email", res.ServiceAccountEmail)
	return nil
}

func resourceGcpServiceAccountCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapGcpServiceAccountFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	ac, err := apiClient.UpdateGcpServiceAccount(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}

	err = mapGcpServiceAccountToData(ac, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGcpServiceAccountRead(ctx, d, meta)
}

func resourceGcpServiceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetGcpServiceAccount(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	err = mapGcpServiceAccountToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGcpServiceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSDCConfig(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGcpServiceAccountImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetGcpServiceAccount(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	err = mapGcpServiceAccountToData(res, d)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
