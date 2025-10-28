package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceSDCConfig() *schema.Resource {
	return &schema.Resource{
		Description:   "Aidbox SDCConfig is a proprietary custom resource used to configure the Aidbox Structured Data Capture (SDC) module.",
		CreateContext: resourceSDCConfigCreate,
		ReadContext:   resourceSDCConfigRead,
		UpdateContext: resourceSDCConfigUpdate,
		DeleteContext: resourceSDCConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSDCConfigImport,
		},
		Schema: resourceFullSchema(resourceSchemaSDCConfig()),
	}
}

func resourceSchemaSDCConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Computer friendly name of the SDC configuration.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"description": {
			Description: "A human-readable description of the SDC configuration.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"default": {
			Description: "Specifies if this is the default configuration for the system or tenant.",
			Type:        schema.TypeBool,
			Optional:    true,
		},
		"storage": {
			Description: "Configuration for storing attachments.",
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket": {
						Description: "The name of the bucket to store attachment files in.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"account": {
						Description: "A reference to the storage account resource.",
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"reference": {
									Description: "The reference to the storage account (e.g. \"GcpServiceAccount/aidbox-rc\").",
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func mapSDCConfigFromData(data *schema.ResourceData) (*aidbox.SDCConfig, error) {
	res := &aidbox.SDCConfig{}
	res.ResourceType = "SDCConfig"
	res.Name = data.Get("name").(string)

	if v, ok := data.GetOk("description"); ok {
		res.Description = v.(string)
	}
	if v, ok := data.GetOk("default"); ok {
		res.Default = v.(bool)
	}

	if v, ok := data.GetOk("storage"); ok {
		storageList := v.([]interface{})
		if len(storageList) > 0 && storageList[0] != nil {
			storageData := storageList[0].(map[string]interface{})
			res.Storage = &aidbox.SDCStorage{}

			if bucket, ok := storageData["bucket"]; ok {
				res.Storage.Bucket = bucket.(string)
			}

			if a, ok := storageData["account"]; ok {
				accountList := a.([]interface{})
				if len(accountList) > 0 && accountList[0] != nil {
					accountData := accountList[0].(map[string]interface{})
					res.Storage.Account = &aidbox.SDCAccount{}

					if ref, ok := accountData["reference"]; ok {
						res.Storage.Account.Reference = ref.(string)
					}
				}
			}
		}
	}

	return res, nil
}

func mapSDCConfigToData(res *aidbox.SDCConfig, data *schema.ResourceData) error {
	data.SetId(res.ID)
	data.Set("name", res.Name)
	data.Set("description", res.Description)
	data.Set("default", res.Default)

	if res.Storage != nil {
		storageData := make(map[string]interface{})
		storageData["bucket"] = res.Storage.Bucket

		if res.Storage.Account != nil {
			accountData := make(map[string]interface{})
			accountData["reference"] = res.Storage.Account.Reference
			storageData["account"] = []interface{}{accountData}
		} else {
			storageData["account"] = []interface{}{}
		}

		data.Set("storage", []interface{}{storageData})
	} else {
		data.Set("storage", []interface{}{})
	}

	return nil
}

func resourceSDCConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSDCConfigFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateSDCConfig(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapSDCConfigToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSDCConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSDCConfig(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	err = mapSDCConfigToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSDCConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSDCConfigFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateSDCConfig(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapSDCConfigToData(ac, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSDCConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSDCConfig(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSDCConfigImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSDCConfig(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	err = mapSDCConfigToData(res, d)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
