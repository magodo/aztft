# AzureRM Terraform Type Finder

`aztft` is a CLI tool (and a library) to query for the AzureRM Terraform Provider resource type based on the input Azure resource ID.

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
|`azurerm_synapse_linked_service`                                 | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Synapse/workspaces/ws1/linkedServices/service1`||
|`azurerm_synapse_role_assignment`                                | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Synapse/workspaces/ws1/roleAssignments/role1`||
|`azurerm_storage_account_queue_properties`                       | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/queueServices/default`||
|`azurerm_storage_account_static_website`                         | `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/account1/staticWebsites/default`||

### Property-like Resources

|Resource Type|Pesudo Resource ID|Comment|
|-|-|-|
|`azurerm_nat_gateway_public_ip_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/natGateways/gw1/publicIPAddresses/<base64 id of azurerm_public_ip>`||
|`azurerm_nat_gateway_public_ip_prefix_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/natGateways/gw1/publicIPPrefixes/<base64 id of azurerm_public_ip_prefix>`||
|`azurerm_network_interface_application_gateway_backend_address_pool_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/ipConfigurations/cfg1/applicationGatewayBackendAddressPools/<base64 of azurerm_application_gateway.example.backend_address_pool.n.id>`||
|`azurerm_network_interface_application_security_group_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/ipConfigurations/cfg1/applicationSecurityGroups/<base64 id of azurerm_application_security_group>`||
|`azurerm_network_interface_backend_address_pool_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/ipConfigurations/cfg1/loadBalancerBackendAddressPools/<base64 id of azurerm_lb_backend_address_pool>`||
|`azurerm_network_interface_nat_rule_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/ipConfigurations/cfg1/loadBalancerInboundNatRules/<base64 id of azurerm_lb_nat_rule>`||
|`azurerm_network_interface_security_group_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1/networkSecurityGruops/<base64 id of azurerm_network_security_group>`||
|`azurerm_subnet_route_table_association`|`/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/subnets/subnet1|routeTables/<base64 id of azurerm_route_table>`||
|`azurerm_subnet_network_security_group_association`|`/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/subnets/subnet1|networkSecurityGroups/<base64 id of azurerm_network_security_group>`||
|`azurerm_subnet_nat_gateway_association`|`/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/subnets/subnet1|natGateways/<base64 id of azurerm_nat_gateway>`||
|`azurerm_virtual_desktop_workspace_application_group_association`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.DesktopVirtualization/workspaces/wsp1/applicationGroups/<base64 id of azurerm_virtual_desktop_application_group>`||
|`azurerm_virtual_machine_data_disk_attachment`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1/dataDisks/disk1`||
|`azurerm_iothub_endpoint_cosmosdb_account`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Devices/iotHubs/hub1/endpointsCosmosdbAccount/ep1`||
|`azurerm_iothub_endpoint_eventhub`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Devices/iotHubs/hub1/endpointsEventhub/ep1`||
|`azurerm_iothub_endpoint_servicebus_queue`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Devices/iotHubs/hub1/endpointsServicebusQueue/ep1`||
|`azurerm_iothub_endpoint_servicebus_topic`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Devices/iotHubs/hub1/endpointsServicebusTopic/ep1`||
|`azurerm_iothub_endpoint_storage_container`| `/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Devices/iotHubs/hub1/endpointsStorageContainer/ep1`||
