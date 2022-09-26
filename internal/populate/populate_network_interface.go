package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateNetworkInterface(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkInterfacesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Interface.Properties
	if props == nil {
		return nil, nil
	}
	nsgProp := props.NetworkSecurityGroup
	if nsgProp == nil {
		return nil, nil
	}

	nsgId := nsgProp.ID
	if nsgId == nil {
	}

	nsgAzureId, err := armid.ParseResourceId(*nsgId)
	if err != nil {
		return nil, fmt.Errorf("parsing resource id %q: %v", *nsgId, err)
	}
	nsgName := nsgAzureId.Names()[0]

	azureId := id.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "networkSecurityGroups")
	azureId.AttrNames = append(azureId.AttrNames, nsgName)

	return []armid.ResourceId{azureId}, nil
}
