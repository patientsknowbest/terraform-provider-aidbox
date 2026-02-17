package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceQuestionnaireTheme() *schema.Resource {
	return &schema.Resource{
		Description:   "QuestionnaireTheme https://www.health-samurai.io/docs/aidbox/reference/system-resources-reference/sdc-module-resources#questionnairetheme",
		CreateContext: resourceQuestionnaireThemeCreate,
		ReadContext:   resourceQuestionnaireThemeRead,
		UpdateContext: resourceQuestionnaireThemeUpdate,
		DeleteContext: resourceQuestionnaireThemeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceQuestionnaireThemeImport,
		},
		Schema: resourceFullSchema(resourceSchemaQuestionnaireTheme()),
	}
}

func resourceSchemaQuestionnaireTheme() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aidbox_id": {
			Description: "The Aidbox ID of the questionnaire theme",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"theme_name": {
			Description: "Name of the theme",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"design_system": {
			Description: "Design system of the theme",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func mapQuestionnaireThemeToData(res *aidbox.QuestionnaireTheme, data *schema.ResourceData) error {
	data.SetId(res.ID)
	if err := data.Set("aidbox_id", res.ID); err != nil {
		return err
	}
	if err := data.Set("theme_name", res.ThemeName); err != nil {
		return err
	}
	if err := data.Set("design_system", res.DesignSystem); err != nil {
		return err
	}
	return nil
}

func mapQuestionnaireThemeFromData(d *schema.ResourceData) *aidbox.QuestionnaireTheme {
	res := &aidbox.QuestionnaireTheme{}
	if v, ok := d.GetOk("aidbox_id"); ok {
		res.ID = v.(string)
	}
	if v, ok := d.GetOk("theme_name"); ok {
		res.ThemeName = v.(string)
	}
	if v, ok := d.GetOk("design_system"); ok {
		res.DesignSystem = v.(string)
	}
	return res
}

func resourceQuestionnaireThemeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapQuestionnaireThemeFromData(d)
	res, err := apiClient.CreateQuestionnaireTheme(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := mapQuestionnaireThemeToData(res, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceQuestionnaireThemeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetQuestionnaireTheme(ctx, d.Id())
	if err != nil {
		if handleNotFoundError(err, d) {
			return nil
		}
		return diag.FromErr(err)
	}
	if err := mapQuestionnaireThemeToData(res, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceQuestionnaireThemeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	q := mapQuestionnaireThemeFromData(d)
	res, err := apiClient.UpdateQuestionnaireTheme(ctx, q)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := mapQuestionnaireThemeToData(res, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceQuestionnaireThemeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	err := apiClient.DeleteQuestionnaireTheme(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceQuestionnaireThemeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetQuestionnaireTheme(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	if err := mapQuestionnaireThemeToData(res, d); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
