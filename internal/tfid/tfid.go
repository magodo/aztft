package tfid

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

type builderFunc func(*client.ClientBuilder, armid.ResourceId) error

var dynamicBuilders = map[string]builderFunc{
	"azurerm_active_directory_domain_service": buildActiveDirectoryDomainService,
}

func NeedsAPI(item resmap.ARMId2TFMapItem) bool {
	_, ok := dynamicBuilders[item.ResourceType]
	return ok
}

func DynamicBuild(id armid.ResourceId, item resmap.ARMId2TFMapItem) (string, error) {
	builder, ok := dynamicBuilders[item.ResourceType]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %q", item.ResourceType)
	}

	b, err := client.NewClientBuilder()
	if err != nil {
		return "", fmt.Errorf("new API client builder: %v", err)
	}

	if err := builder(b, id); err != nil {
		return "", fmt.Errorf("building id for %s: %v", id, err)
	}

	if item.ImportSpec != "" {
		if err := id.Normalize(item.ImportSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q with import spec %q: %v", id.String(), item.ImportSpec, err)
		}
	}
	return id.String(), nil
}

func StaticBuild(id armid.ResourceId, item resmap.ARMId2TFMapItem) (string, error) {
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
	}

	if item.ImportSpec != "" {
		if err := rid.Normalize(item.ImportSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", id.String(), item.ResourceType, item.ImportSpec, err)
		}
	}
	return id.String(), nil
}
