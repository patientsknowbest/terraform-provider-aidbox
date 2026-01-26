package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceStructureDefinitionOverride() *schema.Resource {
	// Note there's no import functionality here. Since the initial create has to read the original spec, it's impossible
	// to import that as that would be already overwritten. This would be probably possible if we asked the user
	// to store the core SD in an attribute, but that adds extra complexity and handling.
	return &schema.Resource{
		Description: "A specialization of StructureDefinition which allows you to override the default version of" +
			" StructureDefinitions that are specified inside the core FHIR IG used on the server. This means default " +
			"rules of resources can be changed without having the client specify a meta.profile in their request.",
		CreateContext: resourceStructureDefinitionOverrideCreate,
		ReadContext:   resourceStructureDefinitionOverrideRead,
		UpdateContext: resourceStructureDefinitionOverrideUpdate,
		DeleteContext: resourceStructureDefinitionOverrideDelete,
		Schema:        resourceFullSchema(resourceSchemaStructureDefinitionOverride()),
	}
}

func resourceSchemaStructureDefinitionOverride() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Description: "Canonical URL that's unique to this StructureDefinition",
			Type:        schema.TypeString,
			Required:    true,
			// url is essentially the logical id, hence it's impossible to update it
			// if you change it, what you mean is: destroy the changed override, and add a new one under this url
			ForceNew: true,
		},
		"structure_definition_override": {
			Description:           "A customized StructureDefinition, based on the original one from the core FHIR spec",
			Type:                  schema.TypeString,
			Required:              true,
			Sensitive:             true,
			DiffSuppressOnRefresh: true,
			DiffSuppressFunc:      jsonDiffSuppressFunc,
		},
		"original_structure_definition": {
			Description: "Backup of the original StructureDefinition, which will be restored upon deleting the override",
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}

func resourceStructureDefinitionOverrideCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)

	// look up the original, FHIR spec StructureDefinition by the canonical URL
	// this is required to create a backup in tf state so we can restore it if we want to delete the
	// customized one - we must never delete  the original as it would break FHIR functionality
	canonicalUrl := d.Get("url").(string)

	originalSD, err := apiClient.GetStructureDefinitionByUrl(ctx, canonicalUrl)

	if err != nil {
		return diag.FromErr(err)
	}

	originalSDBytes, err := json.Marshal(originalSD)
	d.Set("original_structure_definition", string(originalSDBytes))

	overrideSDString := d.Get("structure_definition_override").(string)

	overrideSD := map[string]interface{}{}
	err = json.Unmarshal([]byte(overrideSDString), &overrideSD)
	if err != nil {
		return diag.FromErr(err)
	}

	// now update the SD to our customized version
	updatedSD, err := apiClient.UpdateStructureDefinitionByUrl(ctx, &overrideSD, canonicalUrl)
	if err != nil {
		return diag.FromErr(err)
	}
	var updatedUrl = (*updatedSD)["url"].(string)
	if updatedUrl != canonicalUrl {
		return diag.FromErr(fmt.Errorf("canonical url of resource unexpectedly changed after update, %s was set on the resource but server responded with %s", canonicalUrl, updatedUrl))
	}
	d.Set("url", updatedUrl)
	// throw away the id we didn't know upfront, it just adds unnecessary complexity here when comparing states
	delete(*updatedSD, "id")
	// terraform mandates that we have an id though, so use the url for that
	d.SetId(updatedUrl)

	updatedSDBytes, err := json.Marshal(updatedSD)
	updatedSDString := string(updatedSDBytes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("structure_definition_override", updatedSDString)

	return nil
}

func resourceStructureDefinitionOverrideRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)

	canonicalUrl := d.Get("url").(string)
	overrideSD, err := apiClient.GetStructureDefinitionByUrl(ctx, canonicalUrl)
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}

	var overrideSDUrl = (*overrideSD)["url"].(string)
	if overrideSDUrl != canonicalUrl {
		return diag.FromErr(fmt.Errorf("canonical url of resource unexpectedly changed during state refresh, %s was set on the resource but server responded with %s", canonicalUrl, overrideSDUrl))
	}
	d.Set("url", overrideSDUrl)
	// throw away the id we didn't know upfront, it just adds unnecessary complexity here when comparing states
	delete(*overrideSD, "id")
	// terraform mandates that we have an id though, so use the url for that
	d.SetId(overrideSDUrl)

	overrideSDBytes, err := json.Marshal(overrideSD)
	overrideSDString := string(overrideSDBytes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("structure_definition_override", overrideSDString)

	return nil
}

func resourceStructureDefinitionOverrideUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)

	canonicalUrl := d.Get("url").(string)
	overrideSDString := d.Get("structure_definition_override").(string)

	overrideSD := map[string]interface{}{}
	err := json.Unmarshal([]byte(overrideSDString), &overrideSD)
	if err != nil {
		return diag.FromErr(err)
	}

	// update the SD to our new customized version
	updatedSD, err := apiClient.UpdateStructureDefinitionByUrl(ctx, &overrideSD, canonicalUrl)
	if err != nil {
		return diag.FromErr(err)
	}
	// throw away the id we didn't know upfront, it just adds unnecessary complexity here when comparing states
	delete(*updatedSD, "id")

	updatedSDBytes, err := json.Marshal(updatedSD)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("structure_definition_override", string(updatedSDBytes))

	return nil
}

func resourceStructureDefinitionOverrideDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)

	// no such thing as deleting from the core fhir spec, this just means we restore the original spec
	canonicalUrl := d.Get("url").(string)
	originalSDString := d.Get("original_structure_definition").(string)

	originalSD := map[string]interface{}{}
	err := json.Unmarshal([]byte(originalSDString), &originalSD)
	if err != nil {
		return diag.FromErr(err)
	}

	// restore the SD to the spec version
	if _, err := apiClient.UpdateStructureDefinitionByUrl(ctx, &originalSD, canonicalUrl); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
