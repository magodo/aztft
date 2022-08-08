package aztft

import (
	"fmt"
	"sort"
	"strings"

	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/resolve"
	"github.com/magodo/aztft/internal/tfid"

	"github.com/magodo/armid"
)

// QueryType queries a given ARM resource ID and returns a list of potential matched Terraform resource type.
// It firstly statically search the known resource mappings. If there are multiple matches and the "allowAPI" is true,
// it will further call Azure API to retrieve additionl information about this resource and return the exact match.
func QueryType(idStr string, allowAPI bool) ([]string, error) {
	l, err := query(idStr, allowAPI)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(l))
	for i := range l {
		out[i] = l[i].ResourceType
	}
	return out, nil
}

// QueryId queries a given ARM resource ID and its resource type, returns the matched Terraform resource ID.
func QueryId(idStr string, rt string, allowAPI bool) (string, error) {
	id, err := armid.ParseResourceId(idStr)
	if err != nil {
		return "", fmt.Errorf("parsing id: %v", err)
	}

	resmap.Init()

	m := resmap.TF2ARMIdMap
	_ = m
	item, ok := resmap.TF2ARMIdMap[rt]
	if !ok {
		return "", fmt.Errorf("unknown resource type %q", rt)
	}

	var importSpec string
	if id.ParentScope() == nil {
		// For root scope resource id, the import spec is guaranteed to be only one.
		importSpec = item.ManagementPlane.ImportSpecs[0]
	} else {
		// Otherwise, there might be multiple import specs, which will need to be matched with the scope.
		idscope := id.ParentScope().ScopeString()
		i := -1
		for idx, scope := range item.ManagementPlane.ParentScopes {
			if strings.EqualFold(scope, idscope) {
				i = idx
			}
		}
		if i == -1 {
			return "", fmt.Errorf("id %q doesn't correspond to resource type %q", idStr, rt)
		}
		importSpec = item.ManagementPlane.ImportSpecs[i]
	}

	var spec string
	if tfid.NeedsAPI(rt) {
		if !allowAPI {
			return "", fmt.Errorf("%s needs call Azure API to build the import spec", rt)
		}
		spec, err = tfid.DynamicBuild(id, rt, importSpec)
	} else {
		spec, err = tfid.StaticBuild(id, rt, importSpec)
	}
	if err != nil {
		return "", fmt.Errorf("failed to build id for %s: %v", rt, err)
	}
	return spec, nil
}

func QueryTypeAndId(idStr string, allowAPI bool) ([]string, []string, error) {
	l, err := query(idStr, allowAPI)
	if err != nil {
		return nil, nil, err
	}
	id, _ := armid.ParseResourceId(idStr)

	var outRts, outSpecs []string
	for _, item := range l {
		outRts = append(outRts, item.ResourceType)
		var spec string
		if tfid.NeedsAPI(item.ResourceType) {
			if !allowAPI {
				return nil, nil, fmt.Errorf("%s needs call Azure API to build the import spec", item.ResourceType)
			}
			spec, err = tfid.DynamicBuild(id, item.ResourceType, item.ImportSpec)
		} else {
			spec, err = tfid.StaticBuild(id, item.ResourceType, item.ImportSpec)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build import spec for %s: %v", item.ResourceType, err)
		}
		outSpecs = append(outSpecs, spec)
	}
	return outRts, outSpecs, nil
}

func query(idStr string, allowAPI bool) ([]resmap.ARMId2TFMapItem, error) {
	id, err := armid.ParseResourceId(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id: %v", err)
	}
	k1 := strings.ToUpper(id.RouteScopeString())

	resmap.Init()

	b, ok := resmap.ARMId2TFMap[k1]
	if !ok {
		return nil, nil
	}

	var k2 string
	if id.ParentScope() != nil {
		k2 = strings.ToUpper(id.ParentScope().ScopeString())
	}

	l, ok := b[k2]
	if !ok {
		l, ok = b[strings.ToUpper(resmap.ScopeAny)]
		if !ok {
			return nil, nil
		}
	}

	if len(l) > 1 && allowAPI {
		rt, err := resolve.Resolve(id)
		if err != nil {
			return nil, err
		}
		for _, item := range l {
			if item.ResourceType == rt {
				l = []resmap.ARMId2TFMapItem{item}
				break
			}
		}
		if len(l) > 1 {
			return nil, fmt.Errorf("the ambiguity list doesn't have an item with resource type %q, please open an issue for this", rt)
		}
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].ResourceType < l[j].ResourceType
	})

	return l, nil

}
