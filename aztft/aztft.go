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
		if tfid.NeedsAPI(item) {
			if !allowAPI {
				return nil, nil, fmt.Errorf("%s needs call Azure API to build the import spec", item.ResourceType)
			}
			spec, err = tfid.DynamicBuild(id, item)
		} else {
			spec, err = tfid.StaticBuild(id, item)
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
