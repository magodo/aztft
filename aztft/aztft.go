package aztft

import (
	"fmt"
	"sort"
	"strings"

	"github.com/magodo/aztft/internal/populate"
	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/resolve"
	"github.com/magodo/aztft/internal/tfid"

	"github.com/magodo/armid"
)

type Type struct {
	AzureId armid.ResourceId
	TFType  string
}

// QueryType queries a given ARM resource ID and returns a list of potential matched Terraform resource type.
// It firstly statically search the known resource mappings. If there are multiple matches and the "allowAPI" is true,
// it will further call Azure API to retrieve additionl information about this resource and return the exact match.
// Additionally, if allowAPI is true and this resource maps to multiple TF resources, then multiple Types will be returned.
func QueryType(idStr string, allowAPI bool) (types []Type, exact bool, err error) {
	return queryType(idStr, allowAPI)
}

// QueryId queries a given ARM resource ID and its resource type, returns the matched Terraform resource ID.
func QueryId(idStr string, rt string, allowAPI bool) (string, error) {
	id, err := armid.ParseResourceId(idStr)
	if err != nil {
		return "", fmt.Errorf("parsing id: %v", err)
	}

	return queryId(id, rt, allowAPI)
}

// QueryTypeAndId is similar to QueryType, except it also returns the Terraform resource ID (having same length as the types).
func QueryTypeAndId(idStr string, allowAPI bool) (types []Type, ids []string, exact bool, err error) {
	types, exact, err = queryType(idStr, allowAPI)
	if err != nil {
		return nil, nil, false, err
	}
	for _, t := range types {
		tfid, err := queryId(t.AzureId, t.TFType, allowAPI)
		if err != nil {
			return nil, nil, false, fmt.Errorf("querying id %q as %q: %v", t.AzureId, t.TFType, err)
		}
		ids = append(ids, tfid)
	}
	return types, ids, exact, nil
}

func queryId(id armid.ResourceId, rt string, allowAPI bool) (string, error) {
	var (
		spec string
		err  error
	)
	if tfid.NeedsAPI(rt) {
		if !allowAPI {
			return "", fmt.Errorf("%s needs call Azure API to build the import spec", rt)
		}
		spec, err = tfid.DynamicBuild(id, rt)
	} else {
		spec, err = tfid.StaticBuild(id, rt)
	}
	if err != nil {
		return "", fmt.Errorf("failed to build id for %s: %v", rt, err)
	}
	return spec, nil
}

func getARMId2TFMapItems(id armid.ResourceId) []resmap.ARMId2TFMapItem {
	resmap.Init()
	k1 := strings.ToUpper(id.RouteScopeString())
	b, ok := resmap.ARMId2TFMap[k1]
	if !ok {
		return nil
	}

	var k2 string
	if id.ParentScope() != nil {
		k2 = strings.ToUpper(id.ParentScope().ScopeString())
	}

	l, ok := b[k2]
	if !ok {
		l, ok = b[strings.ToUpper(resmap.ScopeAny)]
		if !ok {
			return nil
		}
	}
	return l
}

func queryType(idStr string, allowAPI bool) ([]Type, bool, error) {
	id, err := armid.ParseResourceId(idStr)
	if err != nil {
		return nil, false, fmt.Errorf("invalid resource id: %v", err)
	}

	l := getARMId2TFMapItems(id)
	if len(l) == 0 {
		return nil, false, nil
	}

	var result []Type

	exact := len(l) == 1
	if allowAPI {
		// Resolve ambiguous resources
		if len(l) > 1 {
			rt, err := resolve.Resolve(id)
			if err != nil {
				return nil, false, err
			}
			for _, item := range l {
				if item.ResourceType == rt {
					l = []resmap.ARMId2TFMapItem{item}
					break
				}
			}
			if len(l) > 1 {
				return nil, false, fmt.Errorf("the ambiguity list doesn't have an item with resource type %q, please open an issue for this", rt)
			}
		}

		// There must be only one resource type, try to populate any property like resources for it.
		exact = true
		result = []Type{
			{
				AzureId: id,
				TFType:  l[0].ResourceType,
			},
		}

		rt := l[0].ResourceType
		propLikeResIds, err := populate.Populate(id, rt)
		if err != nil {
			return nil, false, fmt.Errorf("populating property-like resources for %s: %v", rt, err)
		}

		for _, propLikeResId := range propLikeResIds {
			tmpl := getARMId2TFMapItems(propLikeResId)
			// The resource id of property like resources are hypothetic "unique" resource id, they should have no ambiguity. Otherwise, it is a bug.
			if len(tmpl) != 1 {
				return nil, false, fmt.Errorf("expect 1 TF resource matched for resource id %q, but got %d. Please open an issue for this", propLikeResId, len(tmpl))
			}
			item := tmpl[0]
			result = append(result, Type{
				AzureId: propLikeResId,
				TFType:  item.ResourceType,
			})
		}
	} else {
		for _, item := range l {
			result = append(result, Type{
				AzureId: id,
				TFType:  item.ResourceType,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].AzureId.String() != result[j].AzureId.String() {
			return result[i].AzureId.String() < result[j].AzureId.String()
		}
		return result[i].TFType < result[j].TFType
	})

	return result, exact, nil
}
