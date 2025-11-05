package provider

import (
	"context"
	"encoding/json"
	"reflect"

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
			Description:           "Configuration for storing attachments, as a raw JSON string.",
			Type:                  schema.TypeString,
			Optional:              true,
			DiffSuppressOnRefresh: true,
			DiffSuppressFunc:      jsonDiffSuppressFunc,
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

	// just parse as an "any json" value without validation
	if v, ok := data.GetOk("storage"); ok {
		rawStorage := v.(string)
		if rawStorage != "" {
			storage := &json.RawMessage{}
			err := json.Unmarshal([]byte(rawStorage), storage)
			if err != nil {
				return nil, err
			}
			res.Storage = storage
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
		storage, err := json.Marshal(res.Storage)
		if err != nil {
			return err
		}
		data.Set("storage", string(storage))
	} else {
		data.Set("storage", "")
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
	q.ID = d.Id()

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

func jsonDiffSuppressFunc(_ string, oldJson string, newJson string, _ *schema.ResourceData) bool {
	if oldJson == "" && newJson != "" {
		return false
	}

	var oldObject interface{}
	err := json.Unmarshal([]byte(oldJson), &oldObject)
	if err != nil {
		panic(err)
	}
	var newObject interface{}
	err = json.Unmarshal([]byte(newJson), &newObject)
	if err != nil {
		panic(err)
	}
	return reflect.DeepEqual(oldObject, newObject)
}
