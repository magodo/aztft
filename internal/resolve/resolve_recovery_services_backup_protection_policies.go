package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveRecoveryServicesBackupProtectionPolicies(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewRecoveryServicesBackupProtectionPoliciesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[0], resourceGroupId.Name, id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.ProtectionPolicyResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy:
		return "azurerm_backup_policy_vm", nil
	case *armrecoveryservicesbackup.AzureFileShareProtectionPolicy:
		return "azurerm_backup_policy_file_share", nil
	default:
		return "", fmt.Errorf("unknown policy type: %T", props)
	}
}
