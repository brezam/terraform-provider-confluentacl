# terraform-provider-confluentacl

Custom terraform provider for creating confluent kafka ACLs (Access Control Lists) without needing to create an 
additional cluster api key

## Important Disclaimer (2024)

I don't have access to Confluent Cloud anymore so this is unsupported, and I have no idea if the APIs are still working.
Use at your own risk.
The provider version `0.1.0` with source 'bruno-zamariola/confluentacl' is functionaly equivalent to version `0.1.1` with 
source 'brezam/confluentacl'.

## Why use this?

The official confluent terraform provider requires a `credentials` block in order to create kafka ACLs. This credential 
block uses a kafka api key in that cluster with ConfluentCloudAdmin permissions, which means there's a need to create
an api key in a cluster in order to create other api keys with specific acls in that same cluster.

This requirement is not present when creating ACLs using Confluent UI or Confluent CLI, so this provider was made in order
to be simpler like that as well

## Installation

```hcl
terraform {
  required_providers {
    confluentacl = {
      version = "0.1.1"
      source  = "brezam/confluentacl"
    }
  }
}

provider "confluentacl" {
  confluent_cloud_api_key    = "xxx" // or use environment variable CONFLUENT_CLOUD_API_KEY
  confluent_cloud_api_secret = "xxx" // or use environment variable CONFLUENT_CLOUD_API_SECRET
}
```

## Example

```hcl
resource "confluentacl_acl" "default" {
  rest_endpoint        = "https://XXXXXXXXXX.eastus2.azure.confluent.cloud"
  cluster_id           = "lkc-123abc"
  service_account_name = "my-service-account"

  resource_type = "TOPIC"
  resource_name = "test"
  pattern_type  = "PREFIXED"
  host          = "*"
  operation     = "READ"
  permission    = "ALLOW"
}
```

The resource only requires `cluster_id`, `rest_endpoint`, `service_account_name` and acl configuration (`resource_type`, `resource_name`, ...). 

The `rest_endpoint` and `cluster_id` can easily be obtained from a `confluent_kafka_cluster` resource or data source from 
the official confluent provider.

See [/examples/acl-without-kafka-api-key/main.tf](/examples/acl-without-kafka-api-key/main.tf) for a more detailed usage


## Resources/Data sources

### [Resource] confluentacl_acl

Creates kafka acls. Similar to the official `confluent_kafka_acl`, but without requiring additional kafka cluster credentials

### [Resource] confluentacl_api_key

Creates api keys using a service account name as opposed to service account id. I'd consider it obsolete as you can achieve
the same result by using the official service account data source and api key resource

### [Data source] confluentacl_schema_registry

Loads schema registry id and schema registry url from an environment id. 
Obsolete as now the official provider supports managing schema registry id as well.
