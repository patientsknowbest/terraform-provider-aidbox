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

func TestAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.submitted_responses", "trigger.0.resource", "QuestionnaireResponse"),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.submitted_responses", "trigger.0.fhir_path_criteria", "%current.status = 'completed'"),
					resource.TestCheckResourceAttrWith("aidbox_aidbox_subscription_topic.submitted_responses", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria_updateCriteria,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_aidbox_subscription_topic.submitted_responses", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.submitted_responses", "trigger.0.fhir_path_criteria", "%current.status = 'completed' or %current.status = 'amended'"),
				),
			},
			{
				Config: testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria_removeCriteria,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_aidbox_subscription_topic.submitted_responses", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_aidbox_subscription_topic.submitted_responses", "trigger.0.fhir_path_criteria", ""),
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
const testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria = `
resource "aidbox_aidbox_subscription_topic" "submitted_responses" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/submitted-responses"
  trigger {
    resource           = "QuestionnaireResponse"
    fhir_path_criteria = "%current.status = 'completed'"
  }
}
`
const testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria_updateCriteria = `
resource "aidbox_aidbox_subscription_topic" "submitted_responses" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/submitted-responses"
  trigger {
    resource           = "QuestionnaireResponse"
    fhir_path_criteria = "%current.status = 'completed' or %current.status = 'amended'"
  }
}
`
const testAccAidboxSubscriptionTopic_triggerWithFhirPathCriteria_removeCriteria = `
resource "aidbox_aidbox_subscription_topic" "submitted_responses" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/submitted-responses"
  trigger {
    resource = "QuestionnaireResponse"
  }
}
`
