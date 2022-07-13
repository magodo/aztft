package tfid

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

type builderFunc func(*client.ClientBuilder, armid.ResourceId, string) (string, error)

var dynamicBuilders = map[string]builderFunc{
	"azurerm_active_directory_domain_service": buildActiveDirectoryDomainService,
	"azurerm_storage_object_replication":      buildStorageObjectReplication,
	"azurerm_storage_share":                   buildStorageShare,
	"azurerm_storage_container":               buildStorageContainer,
	"azurerm_storage_queue":                   buildStorageQueue,
	"azurerm_storage_table":                   buildStorageTable,
	"azurerm_key_vault_key":                   buildKeyVaultKey,
	"azurerm_key_vault_secret":                buildKeyVaultSecret,
	"azurerm_key_vault_certificate":           buildKeyVaultCertificate,
}

func NeedsAPI(item resmap.ARMId2TFMapItem) bool {
	_, ok := dynamicBuilders[item.ResourceType]
	return ok
}

func DynamicBuild(id armid.ResourceId, item resmap.ARMId2TFMapItem) (string, error) {
	id = id.Clone()
	builder, ok := dynamicBuilders[item.ResourceType]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %q", item.ResourceType)
	}

	b, err := client.NewClientBuilder()
	if err != nil {
		return "", fmt.Errorf("new API client builder: %v", err)
	}

	return builder(b, id, item.ImportSpec)
}

func StaticBuild(id armid.ResourceId, item resmap.ARMId2TFMapItem) (string, error) {
	id = id.Clone()
	rid, ok := id.(*armid.ScopedResourceId)
	if !ok {
		return id.String(), nil
	}

	switch item.ResourceType {
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
	}

	if item.ImportSpec != "" {
		if err := rid.Normalize(item.ImportSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", id.String(), item.ResourceType, item.ImportSpec, err)
		}
	}
	return id.String(), nil
}
