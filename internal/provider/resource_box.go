package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceBox() *schema.Resource {
	return &schema.Resource{
		Description:   "Box https://docs.aidbox.app/multibox/multibox-box-manager-api.",
		CreateContext: resourceBoxCreate,
		ReadContext:   resourceBoxRead,
		DeleteContext: resourceBoxDelete,
		Importer:      &schema.ResourceImporter{
			StateContext: resourceBoxImport,
		},
		Schema:        resourceSchemaBox(),
	}
}

func mapBoxToData(res *aidbox.Box, data *schema.ResourceData) {
	data.SetId(res.ID)
	data.Set("name", res.ID)
	data.Set("fhir_version", res.FhirVersion)
	data.Set("description", res.Description)
	data.Set("box_url", res.BoxURL)
	var env_vars []interface{}
	for key, value := range res.Env {
		env_vars = append(env_vars, fmt.Sprint(key+"="+value))
	}
	data.Set("env", env_vars)
}

func mapBoxFromData(d *schema.ResourceData) *aidbox.Box {
	res := &aidbox.Box{}
	res.ID = d.Get("name").(string)
	res.Description = d.Get("description").(string)
	res.FhirVersion = d.Get("fhir_version").(string)
	res.BoxURL = d.Get("box_url").(string)
	env := d.Get("env").(*schema.Set).List()
	env_vars := map[string]string{}
	for _, env_var := range env {
		str := fmt.Sprintf("%v", env_var)
		key_value := strings.Split(str, "=")
		env_vars[key_value[0]] = key_value[1]
	}
	res.Env = env_vars
	return res
}

func resourceBoxCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapBoxFromData(d)
	res, err := apiClient.CreateBox(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	mapBoxToData(res, d)
	return nil
}

func resourceBoxRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetBox(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapBoxToData(res, d)
	return nil
}

func resourceBoxDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteBox(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceBoxImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetBox(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	mapBoxToData(res, d)
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaBox() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "name (required): id of the box to create. Must match /[a-z][a-z0-9]{4,}/",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true, // There is no update operation supported right now, so any change forces recreation.
		},
		"description": {
			Description: "Description of box for human users.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"fhir_version": {
			Description: "FHIR version. Value must be from the multibox/versions response.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"env": {
			Description: "object with environment variables in lower-kebab-case (not in UPPER_SNAKE_CASE).",
			Type:        schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			ForceNew: true,
		},
		"box_url": {
			Description: "URL for accessing the box",
			Type:        schema.TypeString,
			Computed:    true,
		},
		//"access_url": {
		//	Type:     schema.TypeString,
		//	Optional: true,
		//	ForceNew: true,
		//},
		//"access_token": {
		//	Type:     schema.TypeString,
		//	Optional: true,
		//	ForceNew: true,
		//},
	}
}
