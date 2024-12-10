package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func baseBoxResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func resourceFullSchema(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	fullSchema := make(map[string]*schema.Schema)
	for k, v := range baseBoxResourceSchema() {
		fullSchema[k] = v
	}
	for k, v := range resourceSchema {
		fullSchema[k] = v
	}
	return fullSchema
}

func mapResourceBaseToData(v *aidbox.ResourceBase, data *schema.ResourceData) {
	if v.ID != "" {
		data.SetId(v.ID)
	}
}

func mapResourceBaseFromData(d *schema.ResourceData) aidbox.ResourceBase {
	res := aidbox.ResourceBase{}
	res.ID = d.Id()
	return res
}
