package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAidboxSubscriptionTopic_triggerPatientEvents(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxSubscriptionTopic_triggerPatientEvents,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.patient_changes", "url", "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"),
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
