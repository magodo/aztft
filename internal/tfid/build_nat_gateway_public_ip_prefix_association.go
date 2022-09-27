package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildNatGatewayPublicIpPrefixAssociation(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
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
	prefixName := id.Names()[1]

	tfNatGwId, err := StaticBuild(natGwId, "azurerm_nat_gateway")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", natGwId, err)
	}

	for _, prefix := range props.PublicIPPrefixes {
		if prefix == nil {
			continue
		}
		if prefix.ID == nil {
			continue
		}
		prefixId, err := armid.ParseResourceId(*prefix.ID)
		if err != nil {
			return "", fmt.Errorf("parsing resource id for %q: %v", *prefix.ID, err)
		}
		if prefixId.Names()[0] != prefixName {
			continue
		}

		tfPrefixId, err := StaticBuild(prefixId, "azurerm_public_ip_prefix")
		if err != nil {
			return "", fmt.Errorf("building resource id for %q: %v", prefixId, err)
		}

		return fmt.Sprintf("%s|%s", tfNatGwId, tfPrefixId), nil
	}

	return "", fmt.Errorf("no nat gateway public ip prefix found by id %q", id)
}
