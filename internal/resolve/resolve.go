package resolve

import (
	"fmt"
	"strings"

	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/resourceid"
)

type resolveFunc func(*client.ClientBuilder, resourceid.ResourceId) (string, error)

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
}

// Resolve resolves a given resource id via Azure API to disambiguate and return a single matched TF resource type.
func Resolve(id resourceid.ResourceId, candidates []resmap.ARMId2TFMapItem) (*resmap.ARMId2TFMapItem, error) {
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
