package provider

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func handleNotFoundError(err error, data *schema.ResourceData) bool {
	if err == aidbox.NotFoundError {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer exists", data.Id())
		data.SetId("")
		return true
	}
	return false
}

func jsonDiffSuppressFunc(_ string, oldJson string, newJson string, _ *schema.ResourceData) bool {
	if oldJson == "" && newJson != "" {
		return false
	}
	if oldJson != "" && newJson == "" {
		return false
	}

	var oldObject interface{}
	err := json.Unmarshal([]byte(oldJson), &oldObject)
	if err != nil {
		panic(err)
	}
	var newObject interface{}
	err = json.Unmarshal([]byte(newJson), &newObject)
	if err != nil {
		panic(err)
	}
	return reflect.DeepEqual(oldObject, newObject)
}
