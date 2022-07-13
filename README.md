# AzureRM Terraform Type Finder

`aztft` is a CLI tool (and a library) to query for the AzureRM Terraform Provider resource type based on the input Azure resource ID.

## Limitation

- `aztft` can only resolves for the main Azure resource's counterpart in Terraform, while those property-like Terraform resources are not handled for now.

- Currently, only Azure management plane resource ID is allowed as input. For Terraform resources that corresponds to Azure resources which are data plane only, we defined following pesudo resource id patterns, which can be recognized by `aztft` as input:

    - `azurerm_key_vault_certificate`                                  : `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1`
	- `azurerm_key_vault_certificate_issuer`                           : `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1/issuers/issuer1`
	- `azurerm_key_vault_managed_storage_account`                      : `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/storage/storage1`
	- `azurerm_key_vault_managed_storage_account_sas_token_definition` : `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/storage/storage1/sas/def1`
	- `azurerm_storage_blob`                                           : `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1/blobs/blob1`
	- `azurerm_storage_data_lake_gen2_filesystem`                      : TBD
	- `azurerm_storage_data_lake_gen2_path`                            : TBD
	- `azurerm_storage_share_directory`                                : TBD
	- `azurerm_storage_table_entity`                                   : TBD
	- `azurerm_storage_share_file`                                     : TBD
	- `azurerm_synapse_role_assignment`                                : TBD
