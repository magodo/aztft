package main

import (
	"fmt"

	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/resolve"
)

func stringInSlice(s string, l []string) bool {
	for _, item := range l {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	resmap.Init()

	// Check whether this resolver's resource types contain deprecated/non-existed resource.
	for k1, m := range resolve.Resolvers {
		for k2, resolver := range m {
			for _, rt := range resolver.ResourceTypes() {
				item, ok := resmap.TF2ARMIdMap[rt]
				if !ok {
					fmt.Printf("non-exist resource type %q for resolver %s in scope of %s\n", rt, k1, k2)
					continue
				}
				if item.IsRemoved {
					fmt.Printf("removed resource type %q for resolver %s in scope of %s\n", rt, k1, k2)
					continue
				}
			}
		}
	}

	for k1, b := range resmap.ARMId2TFMap {
		for k2, l := range b {
			if len(l) > 1 {
				if m, ok := resolve.Resolvers[k1]; ok {
					if resolver, ok := m[k2]; ok {
						// Check whether all the TF candidates are covered by this resolver.
						for _, rt := range l {
							if !stringInSlice(rt.ResourceType, resolver.ResourceTypes()) {
								fmt.Printf("%s in scope of %s has ambiguous resource type %q that isn't covered by that resolver\n", k1, k2, rt)
							}
						}
					}
					continue
				}
				resourceTypes := []string{}
				for _, item := range l {
					resourceTypes = append(resourceTypes, item.ResourceType)
				}
				fmt.Printf("multiple matches found for %s in scope of %s: %v\n", k1, k2, resourceTypes)
			}
		}
	}
}
