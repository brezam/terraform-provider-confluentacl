// // Official confluent provider resources.

data "confluent_environment" "default" {
  display_name = "my-environment-01"
}

// A resource would also work
data "confluent_kafka_cluster" "default" {
  display_name = "my-kafka-cluster-01"
  environment {
    id = data.confluent_environment.default.id
  }
}

// // This provider

resource "confluentacl_acl" "default" {
  service_account_name = "my-service-account"
  rest_endpoint        = data.confluent_kafka_cluster.rest_endpoint
  cluster_id           = data.confluent_kafka_cluster.default.id

  resource_type = "TOPIC"
  resource_name = "test"
  pattern_type  = "PREFIXED"
  host          = "*"
  operation     = "READ"
  permission    = "ALLOW"
}

// While the ACLs are not related to api keys (they're related to service accounts), you can easily see the ACLs 
// by creating an api key and checking it in the Confluent UI page
resource "confluentacl_api_key" "default" {
  service_account_name = "my-service-account"
  environment_id       = data.confluent_environment.default.id
  resource_id          = data.confluent_kafka_cluster.default.id
}
