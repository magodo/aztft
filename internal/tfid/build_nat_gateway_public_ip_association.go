package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildNatGatewayPublicIpAssociation(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkNatGatewaysClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.NatGateway.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	natGwId := id.Parent()
	pipName := id.Names()[1]

	tfNatGwId, err := StaticBuild(natGwId, "azurerm_nat_gateway")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", natGwId, err)
	}

	for _, pip := range props.PublicIPAddresses {
		if pip == nil {
			continue
		}
		if pip.ID == nil {
			continue
		}
		pipId, err := armid.ParseResourceId(*pip.ID)
		if err != nil {
			return "", fmt.Errorf("parsing resource id for %q: %v", *pip.ID, err)
		}
		if pipId.Names()[0] != pipName {
			continue
		}

		tfPipId, err := StaticBuild(pipId, "azurerm_public_ip")
		if err != nil {
			return "", fmt.Errorf("building resource id for %q: %v", pipId, err)
		}

		return fmt.Sprintf("%s|%s", tfNatGwId, tfPipId), nil
	}

	return "", fmt.Errorf("no nat gateway public ip found by id %q", id)
}
