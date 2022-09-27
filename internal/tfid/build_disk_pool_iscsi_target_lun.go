package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildDiskPoolIscsiTargetLun(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStoragePoolIscsiTargetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.IscsiTarget.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	iscsiTargetId := id.Parent()
	diskName := id.Names()[2]

	tfIscsiTargetId, err := StaticBuild(iscsiTargetId, "azurerm_disk_pool_iscsi_target")
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", iscsiTargetId, err)
	}

	for _, lun := range props.Luns {
		if lun == nil {
			continue
		}
		if lun.ManagedDiskAzureResourceID == nil {
			continue
		}
		diskId, err := armid.ParseResourceId(*lun.ManagedDiskAzureResourceID)
		if err != nil {
			return "", fmt.Errorf("parsing resource id for %q: %v", *lun.ManagedDiskAzureResourceID, err)
		}
		if diskId.Names()[0] != diskName {
			continue
		}

		tfDiskId, err := StaticBuild(diskId, "azurerm_managed_disk")
		if err != nil {
			return "", fmt.Errorf("building resource id for %q: %v", diskId, err)
		}

		return fmt.Sprintf("%s/lun|%s", tfIscsiTargetId, tfDiskId), nil
	}

	return "", fmt.Errorf("no disk pool iscsi target lun found by id %q", id)
}
