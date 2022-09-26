package tfid

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildNetworkInterfaceNatRuleAssociation(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
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

	ipConfigId := id.Parent().Parent()
	nicId := ipConfigId.Parent()
	ipConfigName, lbName, natRuleName := id.Names()[1], id.Names()[2], id.Names()[3]

	tfNicId, err := StaticBuild(nicId, "azurerm_network_interface")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", nicId, err)
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
		for _, natRule := range ipConfigProps.LoadBalancerInboundNatRules {
			if natRule.ID == nil {
				continue
			}
			natRuleId, err := armid.ParseResourceId(*natRule.ID)
			if err != nil {
				return "", fmt.Errorf("parsing %q: %v", *natRule.ID, err)
			}
			if !strings.EqualFold(natRuleId.Names()[0], lbName) || !strings.EqualFold(natRuleId.Names()[1], natRuleName) {
				continue
			}

			tfNatRuleId, err := StaticBuild(natRuleId, "azurerm_lb_nat_rule")
			if err != nil {
				return "", fmt.Errorf("building resource id for %q: %v", natRuleId, err)
			}

			return fmt.Sprintf("%s/ipConfigurations/%s|%s", tfNicId, ipConfigName, tfNatRuleId), nil
		}
	}

	return "", fmt.Errorf("no load balancer inbound NAT rule found by id %q", id)
}
