package resmap

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

var (
	//go:embed map.json
	mappingContent []byte

	ARMId2TFMap armId2TFMap

	once sync.Once
)

func Init() {
	once.Do(func() {
		var m TF2ARMIdMap
		if err := json.Unmarshal(mappingContent, &m); err != nil {
			panic(err.Error())
		}
		var err error
		if ARMId2TFMap, err = m.toARM2TFMap(); err != nil {
			panic(err.Error())
		}
	})
}

// TF2ARMIdMap maps from TF resource type to the ARM item
type TF2ARMIdMap map[string]TF2ARMIdMapItem

type TF2ARMIdMapItem struct {
	ManagementPlane *MapManagementPlane `json:"management_plane,omitempty"`

	// Indicates whether this TF resource is removed/deprecated
	IsRemoved bool `json:"is_removed"`
}

const ScopeAny string = "any"

type MapManagementPlane struct {
	// ParentScope is the parent scope in its scope string literal form.
	// Specially:
	// - This is empty for root scope resource ids
	// - A special string "any" means any scope
	ParentScopes []string `json:"scopes,omitempty"`
	Provider     string   `json:"provider"`
	Types        []string `json:"types"`

	// ImportSpecs is a list of ScopeString of valid import specs. They are used to normalize the transformed TF reosurce id, for resolving casing differences (as terraform is case sensitive).
	// Each item should correspond to the item in the ParentScopes, representing a valid import spec in that parent scope.
	// Exceptionally, this might be empty given no import spec is available. This maybe because the parent scope is "any", or this is a root scope resource id.
	ImportSpecs []string `json:"import_specs,omitempty"`
}

// armId2TFMap maps from "<provider>/<types>" (routing scope) to "<parent scope string> | any" to the TF item(s)
type armId2TFMap map[string]map[string][]armId2TFMapItem

type armId2TFMapItem struct {
	ResourceType string
	ImportSpec   string
}

func (mps TF2ARMIdMap) toARM2TFMap() (armId2TFMap, error) {
	out := armId2TFMap{}
	for rt, item := range mps {
		if item.IsRemoved {
			continue
		}
		if item.ManagementPlane == nil {
			continue
		}
		mm := item.ManagementPlane
		k1 := buildRoutingScopeKey(mm.Provider, mm.Types, mm.ParentScopes == nil)

		b, ok := out[k1]
		if !ok {
			b = map[string][]armId2TFMapItem{}
			out[k1] = b
		}

		// The id represents a root scope
		if mm.ParentScopes == nil {
			k2 := ""
			b[k2] = append(b[k2], armId2TFMapItem{
				ResourceType: rt,
			})
			continue
		}

		// The id represents a scoped resource
		for i, ps := range mm.ParentScopes {
			k2 := strings.ToUpper(ps)
			item := armId2TFMapItem{
				ResourceType: rt,
			}
			// Not every item has import spec, this might due to multiple reasons, e.g.:
			// - The TF resource id is synthetic
			// - The TF resource id is under any scope
			// - The TF resource id is a data plane URL
			// - etc...
			// For the items without ImportSpec, it needs a special handling to construct the import spec.
			if len(mm.ImportSpecs) > i {
				item.ImportSpec = mm.ImportSpecs[i]
			}
			b[k2] = append(b[k2], item)
		}
	}
	return out, nil
}

func buildRoutingScopeKey(provider string, types []string, isRootScope bool) string {
	if isRootScope && provider == "Microsoft.Resources" {
		return "/" + strings.ToUpper(strings.Join(types, "/"))
	}
	segs := []string{provider}
	segs = append(segs, types...)
	return "/" + strings.ToUpper(strings.Join(segs, "/"))
}
