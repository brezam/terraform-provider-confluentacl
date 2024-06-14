terraform {
  required_version = "1.5.3"

  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "1.48.0"
    }
    confluentacl = {
      source  = "brezam/confluentacl"
      version = "0.1.1"
    }
  }
}

provider "confluent" {}

provider "confluentacl" {}
