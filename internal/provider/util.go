package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"
	"log"
)

func handleNotFoundError(err error, data *schema.ResourceData) bool {
	if err == aidbox.NotFoundError {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer exists", data.Id())
		data.SetId("")
		return true
	}
	return false
}

func boxIdFromData(d *schema.ResourceData) string {
	boxId, ok := d.GetOk("box_id")
	if ok {
		return boxId.(string)
	} else {
		return ""
	}
}
