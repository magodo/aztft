package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateNatGateway(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkNatGatewaysClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.NatGateway.Properties
	if props == nil {
		return nil, nil
	}

	var result []armid.ResourceId

	for _, pip := range props.PublicIPAddresses {
		if pip == nil {
			continue
		}
		if pip.ID == nil {
			continue
		}
		pipId, err := armid.ParseResourceId(*pip.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", *pip.ID, err)
		}
		pipName := pipId.Names()[0]

		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "publicIPAddresses")
		azureId.AttrNames = append(azureId.AttrNames, pipName)

		result = append(result, azureId)
	}

	return result, nil
}
