package resmap

import (
	_ "embed"
	"encoding/json"
	"strings"
)

var (
	//go:embed map.json
	mappingContent []byte

	ARMId2TFMap armId2TFMap
)

func init() {
	var m tf2ARMIdMap
	if err := json.Unmarshal(mappingContent, &m); err != nil {
		panic(err.Error())
	}
	var err error
	if ARMId2TFMap, err = m.toARM2TFMap(); err != nil {
		panic(err.Error())
	}
}

// tf2ARMIdMap maps from TF resource type to the ARM item
type tf2ARMIdMap map[string]tf2ARMIdMapItem

type tf2ARMIdMapItem struct {
	ManagementPlane *mapManagementPlane `json:"management_plane,omitempty"`
}

const ScopeAny string = "any"

type mapManagementPlane struct {
	// ParentScope is the parent scope in its scope string literal form.
	// Specially:
	// - This is empty for root scope resource ids
	// - A special string "any" means any scope
	ParentScopes []string `json:"scopes,omitempty"`
	Provider     string   `json:"provider"`
	Types        []string `json:"types"`
	Formatter    string   `json:"formatter"`
}

// armId2TFMap maps from "<provider>/<types>" (routing scope) to "<parent scope string> | any" to the TF item(s)
type armId2TFMap map[string]map[string][]armId2TFMapItem

type armId2TFMapItem struct {
	ResourceType string
	Formatter    string
}

func BuildRoutingScopeKey(provider string, types []string) string {
	segs := []string{provider}
	segs = append(segs, types...)
	return strings.ToUpper(strings.Join(segs, "/"))
}

func (mps tf2ARMIdMap) toARM2TFMap() (armId2TFMap, error) {
	out := armId2TFMap{}
	for rt, item := range mps {
		if item.ManagementPlane == nil {
			continue
		}
		mm := item.ManagementPlane
		k1 := BuildRoutingScopeKey(mm.Provider, mm.Types)

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
				Formatter:    mm.Formatter,
			})
			continue
		}

		// The id represents a scoped resource
		for _, ps := range mm.ParentScopes {
			k2 := strings.ToUpper(ps)
			b[k2] = append(b[k2], armId2TFMapItem{
				ResourceType: rt,
				Formatter:    mm.Formatter,
			})
		}
	}
	return out, nil
}
