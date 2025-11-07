package provider

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceSearchParameter() *schema.Resource {
	return &schema.Resource{
		Description:        "SearchParameter https://docs.aidbox.app/api-1/fhir-api/search-1/searchparameter",
		CreateContext:      resourceSearchParameterCreate,
		ReadContext:        resourceSearchParameterRead,
		UpdateContext:      resourceSearchParameterUpdate,
		DeleteContext:      resourceSearchParameterDelete,
		DeprecationMessage: "Legacy implementation, use the resource aidbox_fhir_search_parameter instead (requires enabling schema mode in your server https://docs.aidbox.app/modules/profiling-and-validation/fhir-schema-validator/setup-aidbox-with-fhir-schema-validation-engine)",
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
		"module": {
			Description: "Module name",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"type": {
			Description: "Type of search parameter",
			Type:        schema.TypeString,
			Required:    true,
		},
		"expression": {
			Description: "Expression for elements to search. " +
				"Accepts three types: name of element / index / filter by pattern in collection. " +
				"For filter, separator (|) must be used: {\"system\": \"phone\"} => \"system|phone\"",
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				// PathArray in SearchParameter.expression is an array of strings, integers or objects
				// but in TypeList, "the items are all of the same type defined by the Elem property"
				Type: schema.TypeString,
			},
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
	}
}

func mapSearchParameterFromData(data *schema.ResourceData) (*aidbox.SearchParameter, error) {
	res := &aidbox.SearchParameter{}
	res.Name = data.Get("name").(string)
	res.Module = data.Get("module").(string)

	// type
	t, err := aidbox.ParseSearchParameterType(data.Get("type").(string))
	if err != nil {
		log.Panicln(err)
	}
	res.Type = t

	// expression
	rawElements := data.Get("expression").([]interface{})
	var convertedElements []interface{}
	for _, e := range rawElements {
		ex := e.(string)
		if strings.Contains(ex, "|") {
			// object - filter by pattern in collection
			filterElements := strings.Split(ex, "|")
			filter := map[string]interface{}{
				filterElements[0]: filterElements[1],
			}
			convertedElements = append(convertedElements, filter)
		} else if index, err := strconv.Atoi(ex); err == nil {
			// integer - index in collection
			convertedElements = append(convertedElements, index)
		} else {
			// string - name of element
			convertedElements = append(convertedElements, ex)
		}
	}
	var expression [][]interface{}
	res.ExpressionElements = append(expression, convertedElements)

	// resource
	ref := data.Get("reference").([]interface{})[0].(map[string]interface{})
	r := aidbox.Reference{
		ResourceId:   ref["resource_id"].(string),
		ResourceType: ref["resource_type"].(string),
	}
	res.Resource = r

	res.ID = res.Resource.ResourceId + "." + res.Name

	return res, nil
}

func mapSearchParameterToData(res *aidbox.SearchParameter, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("name", res.Name)
	data.Set("module", res.Module)
	data.Set("type", res.Type.ToString())

	// expression
	var expression []interface{}
	for _, e := range res.ExpressionElements[0] {
		stringElem, ok := e.(string)
		if !ok {
			// number is not an integer thanks to json.Unmarshal in parseResource: https://pkg.go.dev/encoding/json#Unmarshal
			// "To unmarshal JSON into an interface value, Unmarshal stores one of these in the interface value:
			// float64, for JSON numbers"
			floatElem, ok := e.(float64)
			if ok {
				expression = append(expression, strconv.FormatFloat(floatElem, 'f', -1, 64))
			} else {
				var mapElem string
				for k, v := range e.(map[string]interface{}) {
					mapElem = k + "|" + v.(string)
				}
				expression = append(expression, mapElem)
			}
		} else {
			expression = append(expression, stringElem)
		}
	}
	data.Set("expression", expression)

	// resource
	var ref []interface{}
	r := map[string]string{
		"resource_id":   res.Resource.ResourceId,
		"resource_type": res.Resource.ResourceType,
	}
	data.Set("reference", append(ref, r))
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
