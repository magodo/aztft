# AzureRM Terraform Type Finder (WIP)

`aztft` is a CLI tool (and a library) to query for the AzureRM Terraform Provider resource type based on the input Azure resource ID.

## Limitation

The knowledge of this mapping is based on the [AzureRM Terraform provider documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest). Currently, We assume the Terraform resource ID is the same as the Azure resource ID (except for casing), and only Azure management plane (instead of data plane) resource ID is allowed as input.
