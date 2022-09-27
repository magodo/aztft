package populate

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

// populateFunc populates the hypothetic azure resource ids that represent the property like resources of the specified resource.
type populateFunc func(*client.ClientBuilder, armid.ResourceId) ([]armid.ResourceId, error)

var populaters = map[string]populateFunc{
	"azurerm_linux_virtual_machine":     populateVirtualMachine,
	"azurerm_windows_virtual_machine":   populateVirtualMachine,
	"azurerm_network_interface":         populateNetworkInterface,
	"azurerm_virtual_desktop_workspace": populateVirtualDesktopWorkspace,
	"azurerm_nat_gateway":               populateNatGateway,
	"azurerm_disk_pool":                 populateDiskPool,
	"azurerm_disk_pool_iscsi_target":    populateDiskPoolIscsiTarget,
}

func NeedsAPI(rt string) bool {
	_, ok := populaters[rt]
	return ok
}

func Populate(id armid.ResourceId, rt string) ([]armid.ResourceId, error) {
	populater, ok := populaters[rt]
	if !ok {
		return nil, nil
	}

	b, err := client.NewClientBuilder()
	if err != nil {
		return nil, fmt.Errorf("new API client builder: %v", err)
	}

	return populater(b, id)
}
