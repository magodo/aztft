package tfid

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type builderFunc func(*client.ClientBuilder, armid.ResourceId, string) (string, error)

var dynamicBuilders = map[string]builderFunc{
	"azurerm_active_directory_domain_service":                        buildActiveDirectoryDomainService,
	"azurerm_storage_object_replication":                             buildStorageObjectReplication,
	"azurerm_storage_share":                                          buildStorageShare,
	"azurerm_storage_container":                                      buildStorageContainer,
	"azurerm_storage_queue":                                          buildStorageQueue,
	"azurerm_storage_table":                                          buildStorageTable,
	"azurerm_key_vault_key":                                          buildKeyVaultKey,
	"azurerm_key_vault_secret":                                       buildKeyVaultSecret,
	"azurerm_key_vault_certificate":                                  buildKeyVaultCertificate,
	"azurerm_key_vault_certificate_issuer":                           buildKeyVaultCertificateIssuer,
	"azurerm_key_vault_managed_storage_account":                      buildKeyVaultStorageAccount,
	"azurerm_key_vault_managed_storage_account_sas_token_definition": buildKeyVaultStorageAccountSasTokenDefinition,
	"azurerm_storage_blob":                                           buildStorageBlob,
	"azurerm_storage_share_directory":                                buildStorageShareDirectory,
	"azurerm_storage_share_file":                                     buildStorageShareFile,
	"azurerm_storage_table_entity":                                   buildStorageTableEntity,
	"azurerm_storage_data_lake_gen2_filesystem":                      buildStorageDfs,
	"azurerm_storage_data_lake_gen2_path":                            buildStorageDfsPath,
}

func NeedsAPI(rt string) bool {
	_, ok := dynamicBuilders[rt]
	return ok
}

func DynamicBuild(id armid.ResourceId, rt, importSpec string) (string, error) {
	id = id.Clone()
	builder, ok := dynamicBuilders[rt]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %q", rt)
	}

	b, err := client.NewClientBuilder()
	if err != nil {
		return "", fmt.Errorf("new API client builder: %v", err)
	}

	return builder(b, id, importSpec)
}

func StaticBuild(id armid.ResourceId, rt, importSpec string) (string, error) {
	id = id.Clone()
	rid, ok := id.(*armid.ScopedResourceId)
	if !ok {
		return id.String(), nil
	}

	switch rt {
	case "azurerm_app_service_slot_virtual_network_swift_connection":
		rid.AttrTypes[2] = "config"
	case "azurerm_iot_time_series_insights_access_policy":
		rid.AttrTypes[1] = "config"
	case "azurerm_synapse_workspace_sql_aad_admin":
		rid.AttrTypes[1] = "sqlAdministrators"
	case "azurerm_monitor_diagnostic_setting":
		// input: <target id>/providers/Microsoft.Insights/diagnosticSettings/setting1
		// tfid : <target id>|setting1
		id = id.ParentScope()
		return id.String() + "|" + rid.Names()[0], nil

	case "azurerm_synapse_role_assignment":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String() + "|" + id.Names()[1], nil
	case "azurerm_postgresql_active_directory_administrator":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	}

	if importSpec != "" {
		if err := rid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", id.String(), rt, importSpec, err)
		}
	}
	return id.String(), nil
}
