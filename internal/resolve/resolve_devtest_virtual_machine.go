package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveDevTestVirtualMachines(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewDevTestVirtualMachinesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.LabVirtualMachine.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	imageRef := props.GalleryImageReference
	if imageRef == nil {
		return "", fmt.Errorf("unexpected nil galleryImageReference in response")
	}

	osType := imageRef.OSType
	if osType == nil {
		return "", fmt.Errorf("unexpected nil galleryImageReference.osType in response")
	}

	switch *osType {
	case "Linux":
		return "azurerm_dev_test_linux_virtual_machine", nil
	case "Windows":
		return "azurerm_dev_test_windows_virtual_machine", nil
	default:
		return "", fmt.Errorf("unknown os type: %s", *osType)
	}
}
