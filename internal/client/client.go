package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights"
)

type ClientBuilder struct {
	credential azcore.TokenCredential
}

func NewClientBuilder() (*ClientBuilder, error) {
	env := "public"
	if v := os.Getenv("ARM_ENVIRONMENT"); v != "" {
		env = v
	}

	var cloudCfg cloud.Configuration
	switch strings.ToLower(env) {
	case "public":
		cloudCfg = cloud.AzurePublic
	case "usgovernment":
		cloudCfg = cloud.AzureGovernment
	case "china":
		cloudCfg = cloud.AzureChina
	default:
		return nil, fmt.Errorf("unknown environment specified: %q", env)
	}

	// Maps the auth related environment variables used in the provider to what azidentity honors.
	os.Setenv("AZURE_TENANT_ID", os.Getenv("ARM_TENANT_ID"))
	os.Setenv("AZURE_CLIENT_ID", os.Getenv("ARM_CLIENT_ID"))
	os.Setenv("AZURE_CLIENT_SECRET", os.Getenv("ARM_CLIENT_SECRET"))
	os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"))

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudCfg,
		},
		TenantID: os.Getenv("ARM_TENANT_ID"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
	}

	return &ClientBuilder{
		credential: cred,
	}, nil
}

var clientOpt = &arm.ClientOptions{
	ClientOptions: policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{
			ApplicationID: "aztft",
			Disabled:      false,
		},
		Logging: policy.LogOptions{
			IncludeBody: true,
		},
	},
}

func (b *ClientBuilder) NewVirtualMachinesClient(subscriptionId string) (*armcompute.VirtualMachinesClient, error) {
	return armcompute.NewVirtualMachinesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewVirtualMachineScaleSetsClient(subscriptionId string) (*armcompute.VirtualMachineScaleSetsClient, error) {
	return armcompute.NewVirtualMachineScaleSetsClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewDevTestVirtualMachinesClient(subscriptionId string) (*armdevtestlabs.VirtualMachinesClient, error) {
	return armdevtestlabs.NewVirtualMachinesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewRecoveryServicesBackupProtectionPoliciesClient(subscriptionId string) (*armrecoveryservicesbackup.ProtectionPoliciesClient, error) {
	return armrecoveryservicesbackup.NewProtectionPoliciesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewDataProtectionBackupPoliciesClient(subscriptionId string) (*armdataprotection.BackupPoliciesClient, error) {
	return armdataprotection.NewBackupPoliciesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewSynapseIntegrationRuntimesClient(subscriptionId string) (*armsynapse.IntegrationRuntimesClient, error) {
	return armsynapse.NewIntegrationRuntimesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewDigitalTwinsEndpointsClient(subscriptionId string) (*armdigitaltwins.EndpointClient, error) {
	return armdigitaltwins.NewEndpointClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryTriggersClient(subscriptionId string) (*armdatafactory.TriggersClient, error) {
	return armdatafactory.NewTriggersClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewAppServiceCertificatesClient(subscriptionId string) (*armappservice.CertificatesClient, error) {
	return armappservice.NewCertificatesClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewKustoDataConnectionsClient(subscriptionId string) (*armkusto.DataConnectionsClient, error) {
	return armkusto.NewDataConnectionsClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewMachineLearningComputeClient(subscriptionId string) (*armmachinelearning.ComputeClient, error) {
	return armmachinelearning.NewComputeClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewTimeSeriesInsightEnvironmentsClient(subscriptionId string) (*armtimeseriesinsights.EnvironmentsClient, error) {
	return armtimeseriesinsights.NewEnvironmentsClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewStorageCacheTargetsClient(subscriptionId string) (*armstoragecache.StorageTargetsClient, error) {
	return armstoragecache.NewStorageTargetsClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}

func (b *ClientBuilder) NewAutomationConnectionClient(subscriptionId string) (*armautomation.ConnectionClient, error) {
	return armautomation.NewConnectionClient(
		subscriptionId,
		b.credential,
		clientOpt,
	)
}
