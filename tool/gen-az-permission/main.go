package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sort"

	"github.com/magodo/aztft/internal/resmap"
)

// This tool is used to generate ARM permissions for every resource type
func main() {
	var action string
	flag.StringVar(&action, "action", "", "the action of the permission")

	flag.Parse()
	if action == "" {
		panic("action is required")
	}

	resmap.Init()

	rts := make([]string, 0)
	provider2TFMap := make(map[string]string)
	for n, rt := range resmap.TF2ARMIdMap {
		if mgmt := rt.ManagementPlane; mgmt != nil {
			rts = append(rts, mgmt.Provider)
			provider2TFMap[mgmt.Provider] = n
		}
	}

	sort.Sort(sort.StringSlice(rts))

	res := make([]string, 0)
	for _, rt := range rts {
		tf := provider2TFMap[rt]
		if mgmt := resmap.TF2ARMIdMap[tf].ManagementPlane; mgmt != nil {
			for _, t := range mgmt.Types {
				res = append(res, fmt.Sprintf("%s/%s/%s", rt, t, action))
			}
		}
	}

	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
