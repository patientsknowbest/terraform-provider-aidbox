package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
	"log"
	"strings"
)

func resourceSearchParameterV2() *schema.Resource {
	return &schema.Resource{
		Description:   "FHIR R4 SearchParameter https://hl7.org/fhir/R4/searchparameter.html",
		CreateContext: resourceSearchParameterV2Create,
		ReadContext:   resourceSearchParameterV2Read,
		UpdateContext: resourceSearchParameterV2Update,
		DeleteContext: resourceSearchParameterV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchParameterV2Import,
		},
		Schema: resourceFullSchema(resourceSchemaSearchParameterV2()),
	}
}

func resourceSchemaSearchParameterV2() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Name of search parameter, used in search query string",
			Type:        schema.TypeString,
			Required:    true,
		},
		"type": {
			Description: "Type of search parameter",
			Type:        schema.TypeString,
			Required:    true,
		},
		"description": {
			Description: "Natural language description of the search parameter",
			Type:        schema.TypeString,
			Required:    true,
		},
		"url": {
			Description: "Canonical identifier for this search parameter, represented as a URI (globally unique)",
			Type:        schema.TypeString,
			Required:    true,
		},
		"code": {
			Description: "Code used in URL",
			Type:        schema.TypeString,
			Required:    true,
		},
		"status": {
			Description: "Value of draft | active | retired | unknown, see https://hl7.org/fhir/R4/valueset-publication-status.html",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "active",
		},
		"base": {
			Description: "The resource type(s) this search parameter applies to",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"expression": {
			Description: "FHIRPath expression that extracts the values, see https://hl7.org/fhir/fhirpath.html",
			Type:        schema.TypeString,
			Required:    true,
		},
	}
}

func mapSearchParameterV2FromData(data *schema.ResourceData) (*aidbox.SearchParameterV2, error) {
	res := &aidbox.SearchParameterV2{}
	res.Name = data.Get("name").(string)
	res.Description = data.Get("description").(string)
	res.Url = data.Get("url").(string)
	res.Code = data.Get("code").(string)
	res.Status = data.Get("status").(string)
	res.Expression = data.Get("expression").(string)
	res.ResourceType = "SearchParameter"

	// base
	rawBase := data.Get("base").([]interface{})
	base := make([]string, len(rawBase))
	for i, v := range rawBase {
		base[i] = v.(string)
	}
	res.Base = base

	// type
	t, err := aidbox.ParseSearchParameterTypeV2(data.Get("type").(string))
	if err != nil {
		log.Panicln(err)
	}
	res.Type = t

	res.ID = strings.Join(res.Base, "-") + "." + res.Name

	return res, nil
}

func mapSearchParameterV2ToData(res *aidbox.SearchParameterV2, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("base", res.Base)
	data.Set("code", res.Code)
	data.Set("description", res.Description)
	data.Set("expression", res.Expression)
	data.Set("name", res.Name)
	data.Set("status", res.Status)
	data.Set("type", res.Type.ToString())
	data.Set("url", res.Url)
}

func resourceSearchParameterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchParameterV2FromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateSearchParameterV2(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterV2ToData(res, d)
	return nil
}

func resourceSearchParameterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearchParameterV2(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapSearchParameterV2ToData(res, d)
	return nil
}

func resourceSearchParameterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchParameterV2FromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateSearchParameterV2(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterV2ToData(ac, d)
	return nil
}

func resourceSearchParameterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSearchParameterV2(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSearchParameterV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearchParameterV2(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapSearchParameterV2ToData(res, d)
	return []*schema.ResourceData{d}, nil
}
