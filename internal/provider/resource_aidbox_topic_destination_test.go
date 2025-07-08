package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAidboxTopicDestinaton_subscribeToPatientEvents(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxTopicDestinatopn_subscribeToPatientEvents,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "topic", "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.0.name", "endpoint"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.0.url", "https://aidbox.requestcatcher.com/patient-webhook"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.1.name", "timeout"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.1.unsigned_int", "30"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.2.name", "maxMessagesInBatch"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.2.unsigned_int", "1"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.3.name", "header"),
					resource.TestCheckResourceAttr("aidbox_aidbox_topic_destination.patient_changes", "parameter.3.string", "User-Agent: Aidbox Server"),
				),
			},
		},
	})
}

const testAccAidboxTopicDestinatopn_subscribeToPatientEvents = `
resource "aidbox_aidbox_subscription_topic" "patient_changes" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"
  trigger {
    resource = "Patient"
  }
}

resource "aidbox_aidbox_topic_destination" "patient_changes" {
  topic = aidbox_aidbox_subscription_topic.patient_changes.url
  kind = "webhook"
  content = "id-only"
  parameter {
    name = "endpoint"
    url = "https://aidbox.requestcatcher.com/patient-webhook"
  }
  parameter {
    name = "timeout"
    unsigned_int = 30
  }
  parameter {
    name = "maxMessagesInBatch"
    unsigned_int = 1
  }
  parameter {
    name = "header"
    string = "User-Agent: Aidbox Server"
  }
  depends_on = [
    aidbox_aidbox_subscription_topic.patient_changes
  ]
}
`
