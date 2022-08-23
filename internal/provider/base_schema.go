package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

//const (
//	aidboxTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
//)

func baseBoxResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"box_id": {
			Description: "ID of box this object lives in",
			Type:        schema.TypeString,
			Default:     "",
			Optional:    true,
			ForceNew:    true, // Changing the box_id always forces replacement
		},
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
	//if v.Meta != nil {
	//	meta := map[string]interface{}{}
	//	if v.Meta.VersionId != "" {
	//		meta["version_id"] = v.Meta.VersionId
	//	}
	//	if v.Meta.CreatedAt != nil {
	//		meta["created_at"] = v.Meta.CreatedAt.Format(aidboxTimeFormat)
	//	}
	//	if v.Meta.LastUpdated != nil {
	//		meta["last_updated"] = v.Meta.LastUpdated.Format(aidboxTimeFormat)
	//	}
	//	data.Set("meta", meta)
	//}
}

func mapResourceBaseFromData(d *schema.ResourceData) aidbox.ResourceBase {
	res := aidbox.ResourceBase{}
	res.ID = d.Id()
	//if v, ok := d.GetOk("meta"); ok {
	//	meta := v.([]interface{})[0].(map[string]interface{}) // Ugly
	//	mm := &aidbox.ResourceBaseMeta{}
	//	if vv, ok := meta["created_at"]; ok {
	//		if vvv, err := time.Parse(aidboxTimeFormat, vv.(string)); err != nil {
	//			mm.CreatedAt = &vvv
	//		}
	//	}
	//	if vv, ok := meta["last_updated"]; ok {
	//		if vvv, err := time.Parse(aidboxTimeFormat, vv.(string)); err != nil {
	//			mm.LastUpdated = &vvv
	//		}
	//	}
	//	if vv, ok := meta["version_id"]; ok {
	//		mm.VersionId = vv.(string)
	//	}
	//	res.Meta = mm
	//}
	return res
}
