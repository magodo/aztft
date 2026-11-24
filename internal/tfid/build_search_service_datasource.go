package tfid

import (
	"context"
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildSearchServiceDataSource(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSearchServicesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	endpoint := props.Endpoint
	if endpoint == nil {
		return "", fmt.Errorf("unexpected nil properties.endpoint in response")
	}
	uri, err := url.Parse(*endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %v", *endpoint, err)
	}
	uri = uri.JoinPath(id.Names()[2])
	return uri.String(), nil
}
