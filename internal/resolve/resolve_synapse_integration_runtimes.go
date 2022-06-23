package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveSynapseIntegrationRuntimes(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewSynapseIntegrationRuntimesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.IntegrationRuntimeResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armsynapse.ManagedIntegrationRuntime:
		return "azurerm_synapse_integration_runtime_azure", nil
	case *armsynapse.SelfHostedIntegrationRuntime:
		return "azurerm_synapse_integration_runtime_self_hosted", nil
	default:
		return "", fmt.Errorf("unknown integration runtime type: %T", props)
	}
}
