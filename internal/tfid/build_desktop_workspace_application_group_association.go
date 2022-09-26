package tfid

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildDesktopWorkspaceApplicationGroupAssociation(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDesktopVirtualizationWorkspacesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Workspace.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	tfWspId, err := StaticBuild(id.Parent(), "azurerm_virtual_desktop_workspace")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", id.Parent(), err)
	}

	applicationGroupName := id.Names()[1]

	for _, applicationGroupRef := range props.ApplicationGroupReferences {
		if applicationGroupRef == nil {
			continue
		}

		applicationGroupId, err := armid.ParseResourceId(*applicationGroupRef)
		if err != nil {
			return "", fmt.Errorf("parsing %q: %v", *applicationGroupRef, err)
		}

		if !strings.EqualFold(applicationGroupName, applicationGroupId.Names()[0]) {
			continue
		}

		tfApplicationGroupId, err := StaticBuild(applicationGroupId, "azurerm_virtual_desktop_application_group")
		if err != nil {
			return "", fmt.Errorf("building resource id for %q: %v", applicationGroupId, err)
		}
		return tfWspId + "|" + tfApplicationGroupId, nil
	}

	return "", fmt.Errorf("no application group found by id %q", id)
}
