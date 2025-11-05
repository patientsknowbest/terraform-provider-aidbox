package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceStructureDefinition() *schema.Resource {
	return &schema.Resource{
		Description:   "FHIR R4 SearchParameter https://hl7.org/fhir/R4/searchparameter.html",
		CreateContext: resourceStructureDefinitionCreate,
		ReadContext:   resourceStructureDefinitionRead,
		UpdateContext: resourceStructureDefinitionUpdate,
		DeleteContext: resourceStructureDefinitionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStructureDefinitionImport,
		},
		Schema: resourceFullSchema(resourceSchemaStructureDefinition()),
	}
}

func resourceSchemaStructureDefinition() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Computer friendly name of the resource",
			Type:        schema.TypeString,
			Required:    true,
		},
		"url": {
			Description: "Canonical URL that's unique to this StructureDefinition",
			Type:        schema.TypeString,
			Required:    true,
		},
		"base_definition": {
			Description: "Definition that this type is constrained/specialized from",
			Type:        schema.TypeString,
			Required:    true,
		},
		"derivation": {
			Description: "Value of specialization | constraint",
			Type:        schema.TypeString,
			Required:    true,
		},
		"abstract": {
			Description: "Whether the structure is abstract",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"type": {
			Description: "The FHIR type defined or constrained by this structure",
			Type:        schema.TypeString,
			Required:    true,
		},
		"status": {
			Description: "Value of draft | active | retired | unknown",
			Type:        schema.TypeString,
			Required:    true,
		},
		"kind": {
			Description: "Value of primitive-type | complex-type | resource | logical",
			Type:        schema.TypeString,
			Required:    true,
		},
		"version": {
			Description: "Business version of the structure definition",
			Type:        schema.TypeString,
			Required:    true,
		},
		"differential": {
			Description:           "The value of StructureDefinition.differential expressed as a raw JSON string value",
			Type:                  schema.TypeString,
			Required:              true,
			DiffSuppressOnRefresh: true,
			DiffSuppressFunc:      jsonDiffSuppressFunc,
		},
	}
}

func mapStructureDefinitionFromData(data *schema.ResourceData) (*aidbox.StructureDefinition, error) {
	res := &aidbox.StructureDefinition{}
	res.ResourceType = "StructureDefinition"
	res.Name = data.Get("name").(string)
	res.Url = data.Get("url").(string)
	res.BaseDefinition = data.Get("base_definition").(string)
	res.Derivation = data.Get("derivation").(string)
	res.Abstract = data.Get("abstract").(bool)
	res.Type = data.Get("type").(string)
	res.Status = data.Get("status").(string)
	res.Kind = data.Get("kind").(string)
	res.Version = data.Get("version").(string)

	// just parse as an "any json" value without validation
	rawDifferential := data.Get("differential").(string)
	differential := &json.RawMessage{}
	err := json.Unmarshal([]byte(rawDifferential), differential)
	if err != nil {
		return nil, err
	}
	res.Differential = differential

	return res, nil
}

func mapStructureDefinitionToData(res *aidbox.StructureDefinition, data *schema.ResourceData) error {
	data.SetId(res.ID)

	data.Set("name", res.Name)
	data.Set("url", res.Url)
	data.Set("base_definition", res.BaseDefinition)
	data.Set("derivation", res.Derivation)
	data.Set("abstract", res.Abstract)
	data.Set("type", res.Type)
	data.Set("status", res.Status)
	data.Set("kind", res.Kind)
	data.Set("version", res.Version)

	differential, err := json.Marshal(res.Differential)
	if err != nil {
		return err
	}
	data.Set("differential", string(differential))

	return nil
}

func resourceStructureDefinitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapStructureDefinitionFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := apiClient.CreateStructureDefinition(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapStructureDefinitionToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceStructureDefinitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetStructureDefinition(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	err = mapStructureDefinitionToData(res, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceStructureDefinitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q, err := mapStructureDefinitionFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ac, err := apiClient.UpdateStructureDefinition(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mapStructureDefinitionToData(ac, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceStructureDefinitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteStructureDefinition(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceStructureDefinitionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetStructureDefinition(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	err = mapStructureDefinitionToData(res, d)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
