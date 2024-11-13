terraform {
  required_version = ">=1.0.0, <2.0.0"
  required_providers {
    aidbox = {
      source  = "patientsknowbest/aidbox"
      version = "0.13.5"
    }
  }
}

provider "aidbox" {
  client_id     = "root"
  client_secret = "secret"
  url           = "http://localhost:8888"
}

resource "aidbox_aidbox_subscription_topic" "patient_changes" {
  url = "https://fhir.yourcompany.com/subscriptiontopic/patient-changes"
  trigger {
    resource = "Patient"
  }
}

resource "aidbox_aidbox_topic_destination" "patient_changes" {
  topic = aidbox_aidbox_subscription_topic.patient_changes.url
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