package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
	"log"
	"strconv"
	"strings"
)

func resourceSearchParameter() *schema.Resource {
	return &schema.Resource{
		Description:   "SearchParameter https://docs.aidbox.app/api-1/fhir-api/search-1/searchparameter",
		CreateContext: resourceSearchParameterCreate,
		ReadContext:   resourceSearchParameterRead,
		UpdateContext: resourceSearchParameterUpdate,
		DeleteContext: resourceSearchParameterDelete,
		Schema:        resourceFullSchema(resourceSchemaSearchParameter()),
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
			Description: "Searchable elements expression",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"reference": {
			Description: "Reference to resource this search param attached to",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
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
			// due to unmarshal in parseResource, the number becomes a float instead of int
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
	res, err := apiClient.CreateSearchParameter(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterToData(res, d)
	return nil
}

func resourceSearchParameterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetSearchParameter(ctx, d.Id(), boxIdFromData(d))
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
	ac, err := apiClient.UpdateSearchParameter(ctx, q, boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	mapSearchParameterToData(ac, d)
	return nil
}

func resourceSearchParameterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteSearchParameter(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
