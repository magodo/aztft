package importspec

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/resmap"
)

func BuildImportSpec(id armid.ResourceId, item resmap.ARMId2TFMapItem) (string, error) {
	rid, ok := id.(*armid.ScopedResourceId)
	if !ok {
		return id.String(), nil
	}

	// TODO: We should copy the rid here to avoid mutate the input id.

	switch item.ResourceType {
	case "azurerm_app_service_slot_virtual_network_swift_connection":
		rid.AttrTypes[2] = "config"
	case "azurerm_iot_time_series_insights_access_policy":
		rid.AttrTypes[1] = "config"
	case "azurerm_synapse_workspace_sql_aad_admin":
		rid.AttrTypes[1] = "sqlAdministrators"
	}

	if err := rid.Normalize(item.ImportSpec); err != nil {
		return "", fmt.Errorf("normalizing id %q for %q: %v", id.String(), item.ResourceType, err)
	}
	return id.String(), nil
}
