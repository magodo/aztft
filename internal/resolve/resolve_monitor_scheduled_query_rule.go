package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type monitorScheduledQueryRulesResolver struct{}

func (monitorScheduledQueryRulesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_monitor_scheduled_query_rules_alert",
		"azurerm_monitor_scheduled_query_rules_log",

		// This is actually not disambiguated yet, as currently there is only one API version used by the SDK:
		// https://github.com/Azure/azure-sdk-for-go/issues/18960.
		// Once above issue is resolved, we can further disambiguate it from the azurerm_monitor_scheduled_query_rules_alert.
		"azurerm_monitor_scheduled_query_rules_alert_v2",
	}
}

func (monitorScheduledQueryRulesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewMonitorScheduledQueryRulesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.LogSearchRuleResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	action := props.Action
	if action == nil {
		return "", fmt.Errorf("unexpected nil properties.action in response")
	}

	switch action.(type) {
	case *armmonitor.AlertingAction:
		return "azurerm_monitor_scheduled_query_rules_alert", nil
	case *armmonitor.LogToMetricAction:
		return "azurerm_monitor_scheduled_query_rules_log", nil
	default:
		return "", fmt.Errorf("unknown monitor scheduled query rule action type: %T", action)
	}
}
