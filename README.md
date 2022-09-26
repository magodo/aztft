# AzureRM Terraform Type Finder

`aztft` is a CLI tool (and a library) to query for the AzureRM Terraform Provider resource type based on the input Azure resource ID.

## Limitation

`aztft` can only resolves for the main Azure resource's counterpart in Terraform, while those property-like Terraform resources are not handled for now.

## Pesudo Resource ID

In most cases, `aztft` accepts Azure management plane resource ID as input. For other rare cases, some Terraform resources do not correspond to Azure management plane resources, which typically means:

1. The resources are data plane only
2. The resources are property-like

For these resources, as they don't have a management plane resource ID, we defined the "pesudo" resource ID for them:

### Data Plane Only Resources

|Resource Type|Pesudo Resource ID|Comment|
|-|-|-|
|`azurerm_key_vault_certificate`                                  | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1`||
|`azurerm_key_vault_certificate_issuer`                           | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/certificates/cert1/issuers/issuer1`||
|`azurerm_key_vault_managed_storage_account`                      | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/storage/storage1`||
|`azurerm_key_vault_managed_storage_account_sas_token_definition` | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1/storage/storage1/sas/def1`||
|`azurerm_storage_blob`                                           | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1/blobs/blob1`||
|`azurerm_storage_data_lake_gen2_filesystem`                      | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/dfs/dfs1`||
|`azurerm_storage_data_lake_gen2_path`                            | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/dfs/dfs1/paths/path1`|For path that is more than one level, use `:` as separator. E.g. `path1` can be `dir1:dir2`|
|`azurerm_storage_share_directory`                                | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1/directories/path1`|For path that is more than one level, use `:` as separator. E.g. `path1` can be `dir1:dir2`|
|`azurerm_storage_share_file`                                     | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/share1/files/path1`|Note: For path that is more than one level, use `:` as separator. E.g. `path1` can be `dir1:file1`|
|`azurerm_storage_table_entity`                                   | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/tableServices/default/tables/table1/partitionKeys/pk1/rowkeys/rk1`||
|`azurerm_synapse_linked_service`                                | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Synapse/workspaces/ws1/linkedServices/service1`||
|`azurerm_synapse_role_assignment`                                | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Synapse/workspaces/ws1/roleAssignments/role1`||

### Property-like Resources

|Resource Type|Pesudo Resource ID|Comment|
|-|-|-|
|`azurerm_virtual_machine_data_disk_attachment`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1/dataDisks/disk1`||
|`azurerm_network_interface_security_group_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/networkSecurityGruops/group1`||
|`azurerm_network_interface_application_gateway_backend_address_pool_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/ipConfigurations/cfg1/backendAddressPools/pool1`||
