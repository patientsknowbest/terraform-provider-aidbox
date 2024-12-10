package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
	"log"
	"strings"
)

func resourceSearchParameter() *schema.Resource {
	return &schema.Resource{
		Description:   "SearchParameter https://docs.aidbox.app/api-1/fhir-api/search-1/searchparameter",
		CreateContext: resourceSearchParameterCreate,
		ReadContext:   resourceSearchParameterRead,
		UpdateContext: resourceSearchParameterUpdate,
		DeleteContext: resourceSearchParameterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchParameterImport,
		},
		Schema: resourceFullSchema(resourceSchemaSearchParameter()),
	}
}

func resourceSchemaSearchParameter() map[string]*schema.Schema {
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

func mapSearchParameterFromData(data *schema.ResourceData) (*aidbox.SearchParameter, error) {
	res := &aidbox.SearchParameter{}
	res.Name = data.Get("name").(string)
	res.Description = data.Get("description").(string)
	res.Url = data.Get("url").(string)
	res.Code = data.Get("code").(string)
	res.Status = data.Get("status").(string)
	res.Expression = data.Get("expression").(string)

	// base
	rawBase := data.Get("base").([]interface{})
	base := make([]string, len(rawBase))
	for i, v := range rawBase {
		base[i] = v.(string)
	}
	res.Base = base

	// type
	t, err := aidbox.ParseSearchParameterType(data.Get("type").(string))
	if err != nil {
		log.Panicln(err)
	}
	res.Type = t

	res.ID = strings.Join(res.Base, "-") + "." + res.Name

	return res, nil
}

func mapSearchParameterToData(res *aidbox.SearchParameter, data *schema.ResourceData) {
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

func resourceSearchParameterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchParameterFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateSearchParameter(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterToData(res, d)
	return nil
}

func resourceSearchParameterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearchParameter(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapSearchParameterToData(res, d)
	return nil
}

func resourceSearchParameterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapSearchParameterFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateSearchParameter(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterToData(ac, d)
	return nil
}

func resourceSearchParameterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSearchParameter(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSearchParameterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearchParameter(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapSearchParameterToData(res, d)
	return []*schema.ResourceData{d}, nil
}
