package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildDiskPoolManagedDiskAttachement(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStoragePoolDiskPoolsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DiskPool.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	diskPoolId := id.Parent()
	diskName := id.Names()[1]

	tfDiskPoolId, err := StaticBuild(diskPoolId, "azurerm_disk_pool")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", diskPoolId, err)
	}

	for _, disk := range props.Disks {
		if disk == nil {
			continue
		}
		if disk.ID == nil {
			continue
		}
		diskId, err := armid.ParseResourceId(*disk.ID)
		if err != nil {
			return "", fmt.Errorf("parsing resource id for %q: %v", *disk.ID, err)
		}
		if diskId.Names()[0] != diskName {
			continue
		}

		tfDiskId, err := StaticBuild(diskId, "azurerm_managed_disk")
		if err != nil {
			return "", fmt.Errorf("building resource id for %q: %v", diskId, err)
		}

		return fmt.Sprintf("%s/managedDisks|%s", tfDiskPoolId, tfDiskId), nil
	}

	return "", fmt.Errorf("no disk pool managed disk found by id %q", id)
}
