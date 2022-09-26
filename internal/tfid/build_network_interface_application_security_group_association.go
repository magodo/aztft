package tfid

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildNetworkInterfaceApplicationSecurityGroupAssociation(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
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
	asgName := id.Names()[2]

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
		for _, asg := range ipConfigProps.ApplicationSecurityGroups {
			if asg.ID == nil {
				continue
			}
			asgId, err := armid.ParseResourceId(*asg.ID)
			if err != nil {
				return "", fmt.Errorf("parsing %q: %v", *asg.ID, err)
			}
			if !strings.EqualFold(asgId.Names()[0], asgName) {
				continue
			}

			tfAsgId, err := StaticBuild(asgId, "azurerm_application_security_group")
			if err != nil {
				return "", fmt.Errorf("building resource id for %q: %v", asgId, err)
			}

			return tfNicId + "|" + tfAsgId, nil
		}
	}

	return "", fmt.Errorf("no application gateway backend address pool found by id %q", id)
}
