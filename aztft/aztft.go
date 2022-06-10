package aztft

import (
	"fmt"
	"github.com/magodo/aztft/internal/resmap"
	"strings"

	"github.com/magodo/aztft/internal/resourceid"
)

func Resolve(idStr string) ([]string, error) {
	id, err := resourceid.ParseResourceId(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id: %v", err)
	}
	k1 := resmap.BuildRoutingScopeKey(id.Provider(), id.Types())

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

	// TODO: For ARM resource ID that has different format than the TF resource id, need transformation.

	var out []string
	for _, item := range l {
		out = append(out, item.ResourceType)
	}
	return out, nil
}
