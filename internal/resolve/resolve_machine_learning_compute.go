package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resourceid"
)

func resolveMachineLearningComputes(b *client.ClientBuilder, id resourceid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*resourceid.ResourceGroup)
	client, err := b.NewMachineLearningComputeClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.ComputeResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armmachinelearning.ComputeInstance:
		return "azurerm_machine_learning_compute_instance", nil
	case *armmachinelearning.SynapseSpark:
		return "azurerm_machine_learning_synapse_spark", nil
	case *armmachinelearning.AmlCompute:
		return "azurerm_machine_learning_compute_cluster", nil
	case *armmachinelearning.AKS:
		return "azurerm_machine_learning_inference_cluster", nil
	default:
		return "", fmt.Errorf("unknown compute resource type %T", props)
	}
}
