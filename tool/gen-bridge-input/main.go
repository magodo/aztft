package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/resmap"
	"github.com/magodo/aztft/internal/tfid"
)

// The property like resources from map.json that have pesudo Azure resource ID defined
// The list is from: tfid.go:StaticBuild()
var propertyLikeRTs = map[string]bool{
	"azurerm_nat_gateway_public_ip_association":                                      true,
	"azurerm_nat_gateway_public_ip_prefix_association":                               true,
	"azurerm_network_interface_application_gateway_backend_address_pool_association": true,
	"azurerm_network_interface_application_security_group_association":               true,
	"azurerm_network_interface_backend_address_pool_association":                     true,
	"azurerm_network_interface_nat_rule_association":                                 true,
	"azurerm_network_interface_security_group_association":                           true,
	"azurerm_virtual_desktop_workspace_application_group_association":                true,
	"azurerm_subnet_nat_gateway_association":                                         true,
	"azurerm_subnet_network_security_group_association":                              true,
	"azurerm_subnet_route_table_association":                                         true,
}

var (
	SubId       = armid.SubscriptionId{Id: "sub1"}
	MgmtGroupId = armid.ManagementGroup{Name: "grp1"}
	RgId        = armid.ResourceGroup{SubscriptionId: "sub1", Name: "rg1"}
	TenantId    = armid.TenantId{}
)

func main() {
	resmap.Init()
	var rts []string
	for rt := range resmap.TF2ARMIdMap {
		rts = append(rts, rt)
	}
	sort.Sort(sort.StringSlice(rts))

	f := hclwrite.NewEmptyFile()
	body := f.Body()

	for _, rt := range rts {
		entry := resmap.TF2ARMIdMap[rt]

		// Resources need dynamically construct its resource ID are mostly data plane resources
		if tfid.NeedsAPI(rt) {
			continue
		}

		if strings.HasPrefix(rt, "fake_") {
			continue
		}

		// The property-like resources are ignored for now, but can be supported in some way in the future
		if propertyLikeRTs[rt] {
			if err := addExecutionBlock(body, rt, "TODO"); err != nil {
				log.Fatal(err)
			}
			continue
		}

		var idstr string

		switch rt {
		case "azurerm_resource_group":
			idstr = RgId.String()
		case "azurerm_management_group":
			idstr = MgmtGroupId.String()
		case "azurerm_subscription":
			idstr = SubId.String()
		default:
			mp := entry.ManagementPlane

			// Construct the scope id if any
			var (
				scopeRaw string
				scopeId  armid.ResourceId
			)
			if len(mp.ParentScopes) == 1 {
				scopeRaw = mp.ParentScopes[0]
				if scopeRaw == resmap.ScopeAny {
					scopeRaw = "/subscriptions/resourceGroups"
				}
			} else if len(mp.ParentScopes) > 1 {
				scopeRaw = mp.ParentScopes[0]
			}
			if scopeRaw != "" {
				scopeId = routeScopeStrToId(scopeRaw)
			}

			// Construct the resource id
			id := routeScopeStrToId("/" + strings.Join(append([]string{mp.Provider}, mp.Types...), "/"))
			if scopeId != nil {
				routeId := id.(*armid.ScopedResourceId)
				routeId.AttrParentScope = scopeId
			}

			var err error
			idstr, err = tfid.StaticBuild(id, rt)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := addExecutionBlock(body, rt, idstr); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%s", f.Bytes())
}

func addExecutionBlock(mainBody *hclwrite.Body, rt string, idstr string) error {
	execBlk := mainBody.AppendNewBlock("execution", []string{rt, "basic"})
	execBody := execBlk.Body()

	expr, err := buildExpression("path", `"${home}/go/bin/terraform-client-import"`)
	if err != nil {
		return err
	}
	execBody.SetAttributeRaw("path", expr.BuildTokens(nil))

	expr, err = buildExpression("args", fmt.Sprintf(`[
	"-path",
	"${home}/go/bin/terraform-provider-azurerm",
	"-type",
	"%s",
	"-id",
	"%s",
]`, rt, idstr))
	if err != nil {
		return err
	}
	execBody.SetAttributeRaw("args", expr.BuildTokens(nil))
	return nil
}

// routeScopeStrToId turns a route scope string to a resource id, with the names part "randomly" generated
func routeScopeStrToId(input string) armid.ResourceId {
	upperInput := strings.ToUpper(input)

	var parentScope armid.ResourceId = &TenantId
	if strings.HasPrefix(upperInput, strings.ToUpper(RgId.ScopeString())) {
		parentScope = &RgId
	} else if strings.HasPrefix(upperInput, strings.ToUpper(SubId.ScopeString())) {
		parentScope = &SubId
	} else if strings.HasPrefix(upperInput, strings.ToUpper(MgmtGroupId.ScopeString())) {
		parentScope = &MgmtGroupId
	}

	left := input[len(parentScope.ScopeString()):]
	if len(left) == 0 {
		return parentScope
	}

	segs := strings.Split(strings.Trim(left, "/"), "/")
	var names []string
	for _, seg := range segs[1:] {
		names = append(names, seg+"1")
	}
	id := armid.ScopedResourceId{
		AttrParentScope: parentScope,
		AttrProvider:    segs[0],
		AttrTypes:       segs[1:],
		AttrNames:       names,
	}
	return &id
}

func buildExpression(name string, value string) (*hclwrite.Expression, error) {
	src := name + " = " + value

	f, diags := hclwrite.ParseConfig([]byte(src), "gen", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse input: %s", diags)
	}

	attr := f.Body().GetAttribute(name)
	if attr == nil {
		return nil, fmt.Errorf("failed to build expression at the get phase. name = %s, value = %s", name, value)
	}

	return attr.Expr(), nil
}
