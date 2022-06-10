package main

import (
	"fmt"
	"github.com/magodo/aztft/internal"
)

func main() {
	for k1, b := range internal.ARMId2TFMap {
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
}
