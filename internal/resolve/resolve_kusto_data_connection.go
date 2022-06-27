package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveKustoDataConnections(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewKustoDataConnectionsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.DataConnectionClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}
	switch model.(type) {
	case *armkusto.EventGridDataConnection:
		return "azurerm_kusto_eventgrid_data_connection", nil
	case *armkusto.EventHubDataConnection:
		return "azurerm_kusto_eventhub_data_connection", nil
	case *armkusto.IotHubDataConnection:
		return "azurerm_kusto_iothub_data_connection", nil
	default:
		return "", fmt.Errorf("unknown data connection type %T", model)
	}
}
