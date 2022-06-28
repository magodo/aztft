package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/armid"
)

func resolveDataFactoryTriggers(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryTriggersClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.TriggerResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdatafactory.BlobEventsTrigger:
		return "azurerm_data_factory_trigger_blob_event", nil
	case *armdatafactory.ScheduleTrigger:
		return "azurerm_data_factory_schedule_trigger", nil
	case *armdatafactory.CustomEventsTrigger:
		return "azurerm_data_factory_trigger_custom_event", nil
	case *armdatafactory.TumblingWindowTrigger:
		return "azurerm_data_factory_trigger_tumbling_window", nil
	default:
		return "", fmt.Errorf("unknown trigger type %T", props)
	}
}
