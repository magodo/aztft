package resolve

import (
	"fmt"
	"strings"

	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/resourceid"
)

type resolveFunc func(*client.ClientBuilder, resourceid.ResourceId) (string, error)

var resolvers = map[string]map[string]resolveFunc{
	"/MICROSOFT.COMPUTE/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": resolveVirtualMachines,
	},
}

// Resolve resolves a given resource id via Azure API to disambiguate and return a single matched TF resource type.
func Resolve(id resourceid.ResourceId) (*resmap.ARMId2TFMapItem, error) {
	resmap.Init()

	ambiguiteMap := map[string]map[string][]resmap.ARMId2TFMapItem{}
	for routeKey, b := range resmap.ARMId2TFMap {
		for parentScopeKey, l := range b {
			if len(l) > 1 {
				m, ok := ambiguiteMap[routeKey]
				if !ok {
					m = map[string][]resmap.ARMId2TFMapItem{}
					ambiguiteMap[routeKey] = m
				}
				m[parentScopeKey] = l
			}
		}
	}

	routeKey := strings.ToUpper(id.RouteScopeString())
	var parentScopeKey string
	if id.ParentScope() != nil {
		parentScopeKey = strings.ToUpper(id.ParentScope().ScopeString())
	}

	var ambiguiteList []resmap.ARMId2TFMapItem
	if m, ok := ambiguiteMap[routeKey]; ok {
		ambiguiteList = m[parentScopeKey]
	}

	// Ensure the input resource id belongs to the known ambiguaties which is derived from converting the static resource mapping from TF2ARM to ARM2TF.
	if len(ambiguiteList) == 0 {
		return nil, fmt.Errorf("%q is not a known ambiguate resource id", id.String())
	}

	// Ensure the API client can be built.
	b, err := client.NewClientBuilder()
	if err != nil {
		return nil, fmt.Errorf("new API client builder: %v", err)
	}

	if m, ok := resolvers[routeKey]; ok {
		if f, ok := m[parentScopeKey]; ok {
			rt, err := f(b, id)
			if err != nil {
				return nil, fmt.Errorf("resolving %q: %v", id, err)
			}
			for _, item := range ambiguiteList {
				if item.ResourceType == rt {
					return &item, nil
				}
			}
			return nil, fmt.Errorf("Program Bug: The ambiguite list doesn't have an item with resource type %q. Please open an issue for this.", rt)
		}
	}

	return nil, fmt.Errorf("the work to disambiguate for %q is in progress", id)
}
