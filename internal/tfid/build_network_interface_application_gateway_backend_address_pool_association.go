package tfid

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkInterfacesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Interface.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	ipConfigId := id.Parent()
	bapName := id.Names()[2]

	tfNicId, err := StaticBuild(id.Parent().Parent(), "azurerm_network_interface")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", id.Parent().Parent(), err)
	}

	for _, ipConfig := range props.IPConfigurations {
		if ipConfig.ID == nil {
			continue
		}
		if !strings.EqualFold(*ipConfig.ID, ipConfigId.String()) {
			continue
		}
		ipConfigProps := ipConfig.Properties
		if ipConfigProps == nil {
			continue
		}
		for _, bap := range ipConfigProps.ApplicationGatewayBackendAddressPools {
			if bap.ID == nil {
				continue
			}
			bapId, err := armid.ParseResourceId(*bap.ID)
			if err != nil {
				return "", fmt.Errorf("parsing %q: %v", *bap.ID, err)
			}
			if !strings.EqualFold(bapId.Names()[1], bapName) {
				continue
			}

			tfAppGwId, err := StaticBuild(bapId.Parent(), "azurerm_application_gateway")
			if err != nil {
				return "", fmt.Errorf("building resource id for %q: %v", bapId.Parent(), err)
			}

			return fmt.Sprintf("%s/ipConfigurations/%s|%s/backendAddressPools/%s", tfNicId, ipConfigId.Names()[1], tfAppGwId, bapId.Names()[1]), nil
		}
	}

	return "", fmt.Errorf("no application gateway backend address pool found by id %q", id)
}
