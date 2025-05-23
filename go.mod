module github.com/magodo/aztft

go 1.23.0

toolchain go1.24.1

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.18.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.9.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement v0.10.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3 v3.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3 v3.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform/v2 v2.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v4 v4.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation v0.9.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2 v2.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices v1.7.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2 v2.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v7 v7.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection/v3 v3.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datashare/armdatashare v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization/v2 v2.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridkubernetes/armhybridkubernetes v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto/v2 v2.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v4 v4.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7 v7.5.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6 v6.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4 v4.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery/v2 v2.4.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts/v2 v2.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2 v2.0.0-beta.4
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.8.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache/v4 v4.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2 v2.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagepool/armstoragepool v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse v0.8.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloads/armworkloads v1.1.0
	github.com/hashicorp/hcl/v2 v2.17.0
	github.com/magodo/armid v0.0.0-20240524082432-7ce06ae46c33
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v2 v2.27.6
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.2 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/zclconf/go-cty v1.13.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
