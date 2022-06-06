package aztft

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	//go:embed mapping/map.json
	mappingContent []byte
)

type TF2ARMIdMapItems map[string]TF2ARMIdMapItem

type TF2ARMIdMapItem struct {
	ManagementPlane *MapManagementPlane `json:"management_plane,omitempty"`
}

type ARMId2TFMapItem struct {
	ResourceType string
	Formatter    string
}

// From "<provider>/<types>" to "<parent scope string> | any" to the matched map item(s)
type ARMId2TFMapItems map[string]map[string][]ARMId2TFMapItem

func (mps TF2ARMIdMapItems) ToARM2TFMapping() (ARMId2TFMapItems, error) {
	out := ARMId2TFMapItems{}
	for rt, item := range mps {
		if item.ManagementPlane == nil {
			continue
		}
		mm := item.ManagementPlane
		segs := []string{mm.Provider}
		segs = append(segs, mm.Types...)
		k1 := strings.ToUpper(strings.Join(segs, "/"))

		b, ok := out[k1]
		if !ok {
			b = map[string][]ARMId2TFMapItem{}
			out[k1] = b
		}

		// The id represents a root scope
		if mm.ParentScopes == nil {
			k2 := ""
			b[k2] = append(b[k2], ARMId2TFMapItem{
				ResourceType: rt,
				Formatter:    mm.Formatter,
			})
			continue
		}

		// The id represents a scoped resource
		for _, ps := range mm.ParentScopes {
			k2 := strings.ToUpper(ps)
			b[k2] = append(b[k2], ARMId2TFMapItem{
				ResourceType: rt,
				Formatter:    mm.Formatter,
			})
		}
	}
	return out, nil
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
	Formatter    string   `json:"formatter"`
}

func Run() error {
	var tf2armMps TF2ARMIdMapItems
	if err := json.Unmarshal(mappingContent, &tf2armMps); err != nil {
		return err
	}
	arm2tfMps, err := tf2armMps.ToARM2TFMapping()
	if err != nil {
		return err
	}
	_ = arm2tfMps
	for k1, b := range arm2tfMps {
		for k2, l := range b {
			if len(l) > 1 {
				resourceTypes := []string{}
				for _, item := range l {
					resourceTypes = append(resourceTypes, item.ResourceType)
				}
				fmt.Printf("multiple matches found for %s in scope of %s: %v\n", k1, k2, resourceTypes)
			}
		}
	}
	return nil
}
