package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func baseResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//"meta": {
		//	Description: "Metadata for the resource",
		//	Type:        schema.TypeList,
		//	Computed:    true,
		//	Elem: &schema.Resource{Schema: map[string]*schema.Schema{
		//		"created_at": {
		//			Description: "The time the resource was created.",
		//			Type:        schema.TypeString,
		//			Computed:    true,
		//		},
		//		"last_updated": {
		//			Description: "The last time the resource was updated.",
		//			Type:        schema.TypeString,
		//			Computed:    true,
		//		},
		//		"version_id": {
		//			Description: "The version of the resource.",
		//			Type:        schema.TypeString,
		//			Computed:    true,
		//		},
		//	}},
		//},
	}
}

func resourceFullSchema(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	fullSchema := make(map[string]*schema.Schema, 0)
	for k, v := range baseResourceSchema() {
		fullSchema[k] = v
	}
	for k, v := range resourceSchema {
		fullSchema[k] = v
	}
	return fullSchema
}
