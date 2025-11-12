package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAidboxSubscriptionTopic_triggerPatientEvents(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxSubscriptionTopic_triggerPatientEvents,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.patient_changes", "url", "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.patient_changes", "trigger.0.resource", "Patient"),
					resource.TestCheckResourceAttrWith("aidbox_aidbox_subscription_topic.patient_changes", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccAidboxSubscriptionTopic_triggerPatientEvents_updateUrl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_aidbox_subscription_topic.patient_changes", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.patient_changes", "url", "https://fhir.yourcompany.com/subscriptiontopic/patient-changes-updated"),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.patient_changes", "trigger.0.resource", "Patient"),
				),
			},
		},
	})
}

const testAccAidboxSubscriptionTopic_triggerPatientEvents = `
resource "aidbox_aidbox_subscription_topic" "patient_changes" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"
  trigger {
    resource = "Patient"
  }
}
`
const testAccAidboxSubscriptionTopic_triggerPatientEvents_updateUrl = `
resource "aidbox_aidbox_subscription_topic" "patient_changes" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/patient-changes-updated"
  trigger {
    resource = "Patient"
  }
}
`
