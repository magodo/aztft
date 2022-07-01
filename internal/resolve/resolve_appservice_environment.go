package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func resolveAppServiceEnvironemnts(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppServiceEnvironmentsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.EnvironmentResource.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	switch *kind {
	case "ASEV2":
		return "azurerm_app_service_environment", nil
	case "ASEV3":
		return "azurerm_app_service_environment_v3", nil
	default:
		return "", fmt.Errorf("unknown kind: %s", *kind)
	}
}
