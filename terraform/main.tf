# this file can be used for testing nats
# run nats & aidbox locally (tutorial here: https://www.health-samurai.io/docs/aidbox/tutorials/subscriptions-tutorials/aidboxtopicsubscription-nats-tutorial)
# once both run and a stream is set up, in place of step 5. "Create AidboxTopicDestination..." you can apply this via tf to set up the topic and destination
resource "aidbox_aidbox_subscription_topic" "subscription_topic" {
  url = "https://fhir.patientsknowbest.com/subscriptiontopic/patient-changes"
  trigger {
    resource = "Patient"
  }
}

resource "aidbox_aidbox_topic_destination" "patient_resource_changes" {
  kind = "nats-jetstream"
  topic = "https://fhir.patientsknowbest.com/subscriptiontopic/patient-changes"

  content = "id-only"

  // NATS server URL
  parameter {
    name = "url"
    string = "nats://host.docker.internal:4222"
  }
  // NATS subject.
  parameter {
    name = "subject"
    string = "changes.patient"
  }

  //// NATS username in Username/Password Authentication.
  //parameter {
  //  name = "username"
  //  string = ""
  //}
  //// NATS password in Username/Password Authentication.
  //parameter {
  //  name = "password"
  //  string = ""
  //}

  // Is not set if 'none' is provided. Otherwise, uses Java SSLContext.getDefault().
  parameter {
    name = "sslContext"
    string = "none"
  }

  // The connection name is used for metrics. It's not strictly necessary, but suggested for production use.
  //parameter {
  //  name = "connectionName"
  //  string = ""
  //}

  // The path to credentials file.
  //parameter {
  //  name = "credentialsFilePath"
  //  string = ""
  //}

  depends_on = [
    aidbox_aidbox_subscription_topic.subscription_topic
  ]
}

terraform {
  required_version = ">=1.2.0, <2.0.0"
  required_providers {
    aidbox = {
      source  = "patientsknowbest/aidbox"
      version = "0.25.1"
    }
  }
}
