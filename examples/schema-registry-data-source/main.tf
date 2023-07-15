data "confluentacl_schema_registry" "default" {
  environment_id = "env-XXXX"
}

output "schema_data" {
  value = data.confluentacl_schema_registry.default
  // Outputs schema registry id and schema registry url
}
