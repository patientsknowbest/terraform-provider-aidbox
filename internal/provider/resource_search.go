package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceSearch() *schema.Resource {
	return &schema.Resource{
		Description:   "Search https://docs.aidbox.app/api/rest-api/aidbox-search#search-resource",
		CreateContext: resourceSearchCreate,
		ReadContext:   resourceSearchRead,
		UpdateContext: resourceSearchUpdate,
		DeleteContext: resourceSearchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchImport,
		},
		Schema: resourceFullSchema(resourceSchemaSearch()),
	}
}

func resourceSchemaSearch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Name of search, used in search query string",
			Type:        schema.TypeString,
			Required:    true,
		},
		"param_parser": {
			Description: "Parser type for the search parameter (allowed values `token`|`reference`)",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"module": {
			Description: "Module name",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"where": {
			Description: "SQL of search",
			Type:        schema.TypeString,
			Required:    true,
		},
		"reference": {
			Description: "Reference to resource this search param attached to",
			// not a TypeMap because "using the Elem block to define specific keys for the map is currently not possible"
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource_id": {
						Description: "The ID of the referenced resource",
						Type:        schema.TypeString,
						Required:    true,
					},
					"resource_type": {
						Description: "The type of the referenced resource",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
			},
		},
		"token_sql": {
			Description: "SQL templates for token parameter handling, used when param_parser = \"token\". See https://docs.aidbox.app/api/rest-api/aidbox-search#token-search",
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    0,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"only_code": {
						Description: "SQL template when only code is provided.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"only_system": {
						Description: "SQL template when only system is provided.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"no_system": {
						Description: "SQL template when no system is provided.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"both": {
						Description: "SQL template when both system and code are provided.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"text": {
						Description: "SQL template for text search.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"text_format": {
						Description: "Format for text search.",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}

func mapSearchFromData(data *schema.ResourceData) (*aidbox.Search, error) {
	res := &aidbox.Search{}
	res.Name = data.Get("name").(string)
	res.ParamParser = data.Get("param_parser").(string)
	res.Module = data.Get("module").(string)
	res.Where = data.Get("where").(string)

	// resource
	ref := data.Get("reference").([]interface{})[0].(map[string]interface{})
	r := aidbox.Reference{
		ResourceId:   ref["resource_id"].(string),
		ResourceType: ref["resource_type"].(string),
	}
	res.Resource = r

	res.ID = res.Resource.ResourceId + "." + res.Name

	if v, ok := data.GetOk("token_sql"); ok {
		ts := v.([]interface{})[0].(map[string]interface{})
		res.TokenSql = &aidbox.TokenSql{
			OnlyCode:   ts["only_code"].(string),
			OnlySystem: ts["only_system"].(string),
			NoSystem:   ts["no_system"].(string),
			Both:       ts["both"].(string),
			Text:       ts["text"].(string),
			TextFormat: ts["text_format"].(string),
		}
	}
	return res, nil
}

func mapSearchToData(res *aidbox.Search, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("name", res.Name)
	data.Set("param_parser", res.ParamParser)
	data.Set("module", res.Module)
	data.Set("where", res.Where)

	// resource
	var ref []interface{}
	r := map[string]string{
		"resource_id":   res.Resource.ResourceId,
		"resource_type": res.Resource.ResourceType,
	}
	data.Set("reference", append(ref, r))

	if res.TokenSql != nil {
		var ts []interface{}
		// ignoring the error -> following the pattern above (where its ignored implicitly)
		_ = data.Set("token_sql", append(ts, map[string]string{
			"only_code":   res.TokenSql.OnlyCode,
			"only_system": res.TokenSql.OnlySystem,
			"no_system":   res.TokenSql.NoSystem,
			"both":        res.TokenSql.Both,
			"text":        res.TokenSql.Text,
			"text_format": res.TokenSql.TextFormat,
		}))
	}
}

func resourceSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateSearch(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchToData(res, d)
	return nil
}

func resourceSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearch(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapSearchToData(res, d)
	return nil
}

func resourceSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateSearch(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchToData(ac, d)
	return nil
}

func resourceSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSearch(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSearchImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearch(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapSearchToData(res, d)
	return []*schema.ResourceData{d}, nil
}
