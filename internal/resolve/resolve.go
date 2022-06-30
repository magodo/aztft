package resolve

import (
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

type resolveFunc func(*client.ClientBuilder, armid.ResourceId) (string, error)

var Resolvers = map[string]map[string]resolveFunc{
	"/MICROSOFT.COMPUTE/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveVirtualMachines,
	},
	"/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveVirtualMachineScaleSets,
	},
	"/MICROSOFT.DEVTESTLAB/LABS/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDevTestVirtualMachines,
	},
	"/MICROSOFT.APIMANAGEMENT/SERVICE/IDENTITYPROVIDERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveApiManagementIdentities,
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveRecoveryServicesBackupProtectionPolicies,
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPFABRICS/PROTECTIONCONTAINERS/PROTECTEDITEMS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveRecoveryServicesBackupProtectedItems,
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataProtectionBackupPolicies,
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPINSTANCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataProtectionBackupInstances,
	},
	"/MICROSOFT.SYNAPSE/WORKSPACES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveSynapseIntegrationRuntimes,
	},
	"/MICROSOFT.DIGITALTWINS/DIGITALTWINSINSTANCES/ENDPOINTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDigitalTwinsEndpoints,
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/TRIGGERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataFactoryTriggers,
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataFactoryDatasets,
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/LINKEDSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataFactoryLinkedServices,
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDataFactoryIntegrationRuntimes,
	},
	"/MICROSOFT.WEB/CERTIFICATES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveAppServiceCertificates,
	},
	"/MICROSOFT.KUSTO/CLUSTERS/DATABASES/DATACONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveKustoDataConnections,
	},
	"/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/COMPUTES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveMachineLearningComputes,
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveTimeSeriesInsightsEnvironment,
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS/EVENTSOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveTimeSeriesInsightsEventSources,
	},
	"/MICROSOFT.STORAGECACHE/CACHES/STORAGETARGETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveStorageCacheTargets,
	},
	"/MICROSOFT.AUTOMATION/AUTOMATIONACCOUNTS/CONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveAutomationConnections,
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveBotServiceBots,
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES/CHANNELS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveBotServiceChannels,
	},
	"/MICROSOFT.SECURITYINSIGHTS/DATACONNECTORS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": resolveSecurityInsightsDataConnectors,
	},
	"/MICROSOFT.SECURITYINSIGHTS/ALERTRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": resolveSecurityInsightsAlertRules,
	},
	"/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/DATASOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveOperationalInsightsDataSources,
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/BINDINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveAppPlatformBindings,
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/DEPLOYMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveAppPlatformDeployments,
	},
	"/MICROSOFT.DATASHARE/ACCOUNTS/SHARES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveDatashareDatasets,
	},
	"/MICROSOFT.HDINSIGHT/CLUSTERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveHDInsightClusters,
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/INPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveStreamAnalyticsInputs,
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/OUTPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveStreamAnalyticsOutputs,
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/FUNCTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveStreamAnalyticsFunctions,
	},
	"/MICROSOFT.INSIGHTS/SCHEDULEDQUERYRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveMonitorScheduledQueryRules,
	},
}

// Resolve resolves a given resource id via Azure API to disambiguate and return a single matched TF resource type.
func Resolve(id armid.ResourceId, candidates []resmap.ARMId2TFMapItem) (*resmap.ARMId2TFMapItem, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates")
	}
	if len(candidates) == 1 {
		return &candidates[0], nil
	}

	routeKey := strings.ToUpper(id.RouteScopeString())
	var parentScopeKey string
	if id.ParentScope() != nil {
		parentScopeKey = strings.ToUpper(id.ParentScope().ScopeString())
	}

	// Ensure the API client can be built.
	b, err := client.NewClientBuilder()
	if err != nil {
		return nil, fmt.Errorf("new API client builder: %v", err)
	}

	if m, ok := Resolvers[routeKey]; ok {
		if f, ok := m[parentScopeKey]; ok {
			rt, err := f(b, id)
			if err != nil {
				return nil, fmt.Errorf("resolving %q: %v", id, err)
			}
			for _, item := range candidates {
				if item.ResourceType == rt {
					return &item, nil
				}
			}
			return nil, fmt.Errorf("Program Bug: The ambiguite list doesn't have an item with resource type %q. Please open an issue for this.", rt)
		}
	}

	return nil, fmt.Errorf("no resolver found for %q (WIP)", id)
}
