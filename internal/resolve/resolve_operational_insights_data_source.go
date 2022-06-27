package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveOperationalInsightsDataSources(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewOperationalInsightsDataSourcesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	switch *kind {
	case armoperationalinsights.DataSourceKindWindowsPerformanceCounter:
		return "azurerm_log_analytics_datasource_windows_performance_counter", nil
	case armoperationalinsights.DataSourceKindWindowsEvent:
		return "azurerm_log_analytics_datasource_windows_event", nil
	default:
		return "", fmt.Errorf("unknown data source kind: %s", *kind)
	}
}
