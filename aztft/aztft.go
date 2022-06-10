package aztft

import (
	"fmt"
	"github.com/magodo/aztft/internal/resmap"
	"sort"
	"strings"

	"github.com/magodo/aztft/internal/resourceid"
	"github.com/magodo/aztft/internal/transformid"
)

func Resolve(idStr string) ([]string, error) {
	id, err := resourceid.ParseResourceId(idStr)
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
		return nil, nil
	}

	// For ARM resource ID that has different format than the TF resource id, need transformation.
	id = transformid.TransformId(id)

	var out []string
	for _, item := range l {
		out = append(out, item.ResourceType)
	}
	sort.Strings(out)
	return out, nil
}
