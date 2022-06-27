package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveSecurityInsightsDataConnectors(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewSecurityInsightsDataConnectorsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.ParentScope().Names()[0], id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.DataConnectorClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}

	switch model.(type) {
	case *armsecurityinsights.MCASDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_cloud_app_security", nil
	case *armsecurityinsights.AADDataConnector:
		return "azurerm_sentinel_data_connector_azure_active_directory", nil
	case *armsecurityinsights.OfficeDataConnector:
		return "azurerm_sentinel_data_connector_office_365", nil
	case *armsecurityinsights.TIDataConnector:
		return "azurerm_sentinel_data_connector_threat_intelligence", nil
	case *armsecurityinsights.AwsS3DataConnector:
		return "azurerm_sentinel_data_connector_aws_s3", nil
	case *armsecurityinsights.AwsCloudTrailDataConnector:
		return "azurerm_sentinel_data_connector_aws_cloud_trail", nil
	case *armsecurityinsights.ASCDataConnector:
		return "azurerm_sentinel_data_connector_azure_security_center", nil
	case *armsecurityinsights.MDATPDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection", nil
	case *armsecurityinsights.AATPDataConnector:
		return "azurerm_sentinel_data_connector_azure_advanced_threat_protection", nil
	default:
		return "", fmt.Errorf("unknown data connector type: %T", model)
	}
}
