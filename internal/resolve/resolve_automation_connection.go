package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveAutomationConnections(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewAutomationConnectionClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Connection.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	connType := props.ConnectionType
	if connType == nil {
		return "", fmt.Errorf("unexpected nil properties.connectionType in response")
	}
	connTypeName := connType.Name
	if connTypeName == nil {
		return "", fmt.Errorf("unexpected nil property.connectionType.name in response")
	}

	switch *connTypeName {
	case "AzureServicePrincipal":
		return "azurerm_automation_connection_service_principal", nil
	case "Azure":
		return "azurerm_automation_connection_certificate", nil
	case "AzureClassicCertificate":
		return "azurerm_automation_connection_classic_certificate", nil
	default:
		return "azurerm_automation_connection", nil
	}
}
