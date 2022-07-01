package resolve

import (
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

type resolver interface {
	Resolve(*client.ClientBuilder, armid.ResourceId) (string, error)
	ResourceTypes() []string
}

var Resolvers = map[string]map[string]resolver{
	"/MICROSOFT.COMPUTE/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualMachinesResolver{},
	},
	"/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualMachineScaleSetsResolver{},
	},
	"/MICROSOFT.DEVTESTLAB/LABS/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": devTestVirtualMachinesResolver{},
	},
	"/MICROSOFT.APIMANAGEMENT/SERVICE/IDENTITYPROVIDERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": apiManagementIdentitiesResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": recoveryServicesBackupProtectionPoliciesResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPFABRICS/PROTECTIONCONTAINERS/PROTECTEDITEMS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": recoveryServicesBackupProtectedItemsResolver{},
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataProtectionBackupPoliciesResolver{},
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPINSTANCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataProtectionBackupInstancesResolver{},
	},
	"/MICROSOFT.SYNAPSE/WORKSPACES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": synapseIntegrationRuntimesResolver{},
	},
	"/MICROSOFT.DIGITALTWINS/DIGITALTWINSINSTANCES/ENDPOINTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": digitalTwinsEndpointsResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/TRIGGERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryTriggersResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryDatasetsResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/LINKEDSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryLinkedServicesResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryIntegrationRuntimesResolver{},
	},
	"/MICROSOFT.KUSTO/CLUSTERS/DATABASES/DATACONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": kustoDataConnectionsResolver{},
	},
	"/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/COMPUTES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": machineLearningComputesResolver{},
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": timeSeriesInsightsEnvironmentResolver{},
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS/EVENTSOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": timeSeriesInsightsEventSourcesResolver{},
	},
	"/MICROSOFT.STORAGECACHE/CACHES/STORAGETARGETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": storageCacheTargetsResolver{},
	},
	"/MICROSOFT.AUTOMATION/AUTOMATIONACCOUNTS/CONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": automationConnectionsResolver{},
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": botServiceBotsResolver{},
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES/CHANNELS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": botServiceChannelsResolver{},
	},
	"/MICROSOFT.SECURITYINSIGHTS/DATACONNECTORS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": securityInsightsDataConnectorsResolver{},
	},
	"/MICROSOFT.SECURITYINSIGHTS/ALERTRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": securityInsightsAlertRulesResolver{},
	},
	"/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/DATASOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": operationalInsightsDataSourcesResolver{},
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/BINDINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appPlatformBindingsResolver{},
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/DEPLOYMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appPlatformDeploymentsResolver{},
	},
	"/MICROSOFT.DATASHARE/ACCOUNTS/SHARES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": datashareDatasetsResolver{},
	},
	"/MICROSOFT.HDINSIGHT/CLUSTERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": hdInsightClustersResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/INPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsInputsResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/OUTPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsOutputsResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/FUNCTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsFunctionsResolver{},
	},
	"/MICROSOFT.INSIGHTS/SCHEDULEDQUERYRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": monitorScheduledQueryRulesResolver{},
	},
	"/MICROSOFT.CDN/PROFILES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": cdnProfilesResolver{},
	},
	"/MICROSOFT.WEB/CERTIFICATES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceCertificatesResolver{},
	},
	"/MICROSOFT.WEB/SITES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSitesResolver{},
	},
	"/MICROSOFT.WEB/SITES/SLOTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSiteSlotsResolver{},
	},
	"/MICROSOFT.WEB/SITES/HYBRIDCONNECTIONNAMESPACES/RELAYS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSiteHybridConnectionsResolver{},
	},
	"/MICROSOFT.WEB/HOSTINGENVIRONMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceEnvironemntsResolver{},
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
		if resolver, ok := m[parentScopeKey]; ok {
			rt, err := resolver.Resolve(b, id)
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
