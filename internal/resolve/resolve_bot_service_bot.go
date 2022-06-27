package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveBotServiceBots(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewBotServiceBotsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Bot.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}

	switch *kind {
	case armbotservice.KindAzurebot:
		return "azurerm_bot_service_azure_bot", nil
	case armbotservice.KindBot:
		return "azurerm_bot_channels_registration", nil
	case armbotservice.KindSdk:
		return "azurerm_bot_web_app", nil
	default:
		return "", fmt.Errorf("unknown bot kind: %s", *kind)
	}
}
