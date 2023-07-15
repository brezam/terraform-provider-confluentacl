---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "confluentacl_acl Resource - terraform-provider-confluentacl"
subcategory: ""
description: |-
  
---

# confluentacl_acl (Resource)

This resources allows creating kafka cluster ACLs (access control lists) without using a preexisting kafka cluster 
api key with credentials. Common usage is as follows:

```terraform
// // Official confluent provider resources.
data "confluent_environment" "default" {
  display_name = "my-environment-01"
}

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
// ACL creation without needing api key/secret in the cluster beforehand
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

- `cluster_id` (String) (Required) ID of the confluent kafka cluster
- `host` (String) (Required) The host for the ACL. Should be set to `*`
- `operation` (String) (Required)  The operation type for the ACL. Possible values: `ALL`, `READ`, `WRITE`, `CREATE`, `DELETE`, `ALTER`, `DESCRIBE`, `CLUSTER_ACTION`, `DESCRIBE_CONFIGS`, `ALTER_CONFIGS`, and `IDEMPOTENT_WRITE`.
- `pattern_type` (String) (Required) The pattern type for the ACL. Possible values: `LITERAL` and `PREFIXED`.
- `permission` (String) (Required) The permission for the ACL. Should be either `DENY` or `ALLOW`.
- `resource_name` (String) (Required) The resource name for the ACL. Must be `kafka-cluster` if `resource_type` equals to `CLUSTER`.
- `resource_type` (String) (Required) The type of the resource. Possible values: `TOPIC`, `GROUP`, `CLUSTER`, `TRANSACTIONAL_ID`, `DELEGATION_TOKEN`.
- `rest_endpoint` (String) (Required) REST endpoint of the kafka cluster
- `service_account_name` (String) (Required) Name of the service account that will be the owner of the ACLs

### Attributes Reference

- `id` (String) The ID of this resource.