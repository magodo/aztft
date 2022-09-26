package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

func buildNetworkInterfaceSecurityGroupAssociation(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
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
	nsgProp := props.NetworkSecurityGroup
	if nsgProp == nil {
		return "", fmt.Errorf("unexpected nil NSG in properties")
	}
	nsgId := nsgProp.ID
	if nsgId == nil {
		return "", fmt.Errorf("unexpected nil NSG Id in properties")
	}

	resmap.Init()

	tfNicId, err := StaticBuild(id.Parent(), "azurerm_network_interface")
	if err != nil {
		return "", fmt.Errorf("building resource id for %s: %v", id.Parent(), err)
	}

	nsgAzureId, err := armid.ParseResourceId(*nsgId)
	if err != nil {
		return "", fmt.Errorf("parsing nsg id %q: %v", *nsgId, err)
	}
	tfNsgId, err := StaticBuild(nsgAzureId, "azurerm_network_security_group")
	if err != nil {
		return "", fmt.Errorf("building resource id for %s: %v", id.Parent(), err)
	}

	return tfNicId + "|" + tfNsgId, nil
}
