terraform {
  required_providers {
    hyperfabric = {
      source = "CiscoDevNet/hyperfabric"
    }
  }
}

provider "hyperfabric" {
  # token = "<MY_HYPERFABRIC_TOKEN>"
}