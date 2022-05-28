package resourceid

import (
	"fmt"
	"strings"
)

type ResourceId interface {
	// ParentScope returns the parent scope of this resource. Normally, scopes are seperated by "/providers/".
	// This is nil if the resource itself is a root scope.
	// E.g.
	// - /subscriptions/0000/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1 	-(parent scope)-> /subscriptions/0000/resourceGroups/rg1
	// - /subscriptions/0000/resourceGroups/rg1 									-(parent scope)-> nil
	ParentScope() ResourceId

	// Parent returns the parent resource. The parent resource belongs to the same provider as the current resource.
	// Nil is returned if the current resource is a root scoped resource, or this is a root scope.
	Parent() ResourceId

	// Provider returns the provider namespace of this resource id.
	// For scoped resource, it is the provider namespace of its routing scope, i.e. the scope of the resource itself.
	// For root scopes, it is the builtin provider namespace: "Microsoft.Resources".
	Provider() string

	// Types returns the resource type array of this resource.
	// For scoped resource, it is the sub-types of its routing scope, i.e. the scope of the resource itself.
	// e.g. ["virtualNetworks", "subnets"] for "Microsoft.Network/virtualNetworks/subnets"
	// For root scopes, it is a builtin type.
	Types() []string

	// Names returns the resource name array of this resource.
	// For scoped resource, it is the names of each sub-type of the Types(), which indicates it always has the same length as the return value of Types().
	// For root scopes, it is nil.
	Names() []string

	// String returns the resource id literal.
	String() string
}

func ParseResourceId(id string) (ResourceId, error) {
	if id == "/" {
		return TenantId{}, nil
	}
	if !strings.HasPrefix(id, "/") {
		return nil, fmt.Errorf(`id should start with "/"`)
	}
	segs := strings.Split(id[1:], "/")

	for idx, seg := range segs {
		if seg == "" {
			return nil, fmt.Errorf(`empty segment found behind %dth "/"`, idx+1)
		}
	}

	var rootScope RootScope = TenantId{}
	if len(segs) >= 4 && segs[0] == "subscriptions" && strings.EqualFold(segs[2], "resourcegroups") {
		rootScope = ResourceGroup{
			SubscriptionId: segs[1],
			Name:           segs[3],
		}
		segs = segs[4:]
	} else if len(segs) >= 2 && segs[0] == "subscriptions" {
		rootScope = SubscriptionId{
			Id: segs[1],
		}
		segs = segs[2:]
	} else if len(segs) >= 4 && segs[0] == "providers" && segs[1] == "Microsoft.Management" && strings.EqualFold(segs[2], "managementgroups") {
		rootScope = ManagementGroup{
			Name: segs[3],
		}
		segs = segs[4:]
	}

	var rid ResourceId = rootScope
	for len(segs) != 0 {
		if segs[0] != "providers" {
			return nil, fmt.Errorf(`scopes should be split by "/providers/"`)
		}
		segs = segs[1:]

		if len(segs) == 0 {
			return nil, fmt.Errorf("missing provider namespace segment")
		}
		provider := segs[0]
		segs = segs[1:]

		var types, names []string

		if len(segs) == 0 || segs[0] == "providers" {
			return nil, fmt.Errorf("missing sub-type type")
		}
		for len(segs) != 0 {
			types = append(types, segs[0])
			segs = segs[1:]

			if len(segs) == 0 {
				return nil, fmt.Errorf("missing sub-type name")
			}
			names = append(names, segs[0])
			segs = segs[1:]

			if len(segs) != 0 && segs[0] == "providers" {
				break
			}
		}
		rid = ScopedResourceId{
			Scope:         rid,
			Namespace:     provider,
			ResourceTypes: types,
			ResourceNames: names,
		}
	}
	return rid, nil
}

// RootScope is a special resource id, that represents a root scope as defined by ARM.
// This is a sealed interface, that has a limited set of implementors.
type RootScope interface {
	ResourceId
	isRootScope()
}

// TenantId represents the tenant scope, which is a pesudo resource id.
type TenantId struct{}

var _ RootScope = TenantId{}

func (TenantId) ParentScope() ResourceId {
	return nil
}

func (TenantId) Parent() ResourceId {
	return nil
}

func (TenantId) Provider() string {
	return "Microsoft.Resources"
}

func (TenantId) Types() []string {
	return []string{"tenants"}
}

func (TenantId) Names() []string {
	return nil
}

func (TenantId) String() string {
	return "/"
}

func (TenantId) isRootScope() {}

// SubscriptionId represents the subscription scope
type SubscriptionId struct {
	// Id is the UUID of this subscription
	Id string
}

var _ RootScope = SubscriptionId{}

func (SubscriptionId) ParentScope() ResourceId {
	return nil
}

func (SubscriptionId) Provider() string {
	return "Microsoft.Resources"
}

func (SubscriptionId) Parent() ResourceId {
	return nil
}

func (SubscriptionId) Types() []string {
	return []string{"subscriptions"}
}

func (SubscriptionId) Names() []string {
	return nil
}

func (id SubscriptionId) String() string {
	return "/subscriptions/" + id.Id
}

func (SubscriptionId) isRootScope() {}

// ResourceGroup represents the resource group scope
type ResourceGroup struct {
	// SubscriptionId is the UUID of the containing subscription
	SubscriptionId string
	// Name is the name of this resource group
	Name string
}

var _ RootScope = ResourceGroup{}

func (ResourceGroup) ParentScope() ResourceId {
	return nil
}

func (ResourceGroup) Parent() ResourceId {
	return nil
}

func (ResourceGroup) Provider() string {
	return "Microsoft.Resources"
}

func (ResourceGroup) Types() []string {
	return []string{"resourceGroups"}
}

func (ResourceGroup) Names() []string {
	return nil
}

func (id ResourceGroup) String() string {
	return "/subscriptions/" + id.SubscriptionId + "/resourceGroups/" + id.Name
}

func (ResourceGroup) isRootScope() {}

// ManagementGroup represents the management group scope
type ManagementGroup struct {
	// Name is the name of this management group
	Name string
}

var _ RootScope = ResourceGroup{}

func (ManagementGroup) ParentScope() ResourceId {
	return nil
}

func (ManagementGroup) Parent() ResourceId {
	return nil
}

func (ManagementGroup) Provider() string {
	return "Microsoft.Management"
}

func (ManagementGroup) Types() []string {
	return []string{"managementGroups"}
}

func (ManagementGroup) Names() []string {
	return nil
}

func (id ManagementGroup) String() string {
	return formatScope(id.Provider(), id.Types(), []string{id.Name})
}

func (ManagementGroup) isRootScope() {}

// ScopedResourceId represents a resource id that is scoped within a root scope or another scoped resource.
var _ ResourceId = ScopedResourceId{}

type ScopedResourceId struct {
	Scope         ResourceId
	Namespace     string
	ResourceTypes []string
	ResourceNames []string
}

func (id ScopedResourceId) ParentScope() ResourceId {
	return id.Scope
}

func (id ScopedResourceId) Parent() ResourceId {
	length := len(id.ResourceTypes)
	if length == 1 {
		return nil
	}
	return ScopedResourceId{
		Scope:         id.Scope,
		Namespace:     id.Namespace,
		ResourceTypes: id.ResourceTypes[0 : length-1],
		ResourceNames: id.ResourceNames[0 : length-1],
	}
}

func (id ScopedResourceId) Provider() string {
	return id.Namespace
}

func (id ScopedResourceId) Types() []string {
	return id.ResourceTypes
}

func (id ScopedResourceId) Names() []string {
	return id.ResourceNames
}

func (id ScopedResourceId) String() string {
	builder := strings.Builder{}
	if _, ok := id.ParentScope().(TenantId); !ok {
		builder.WriteString(id.ParentScope().String())
	}
	builder.WriteString(formatScope(id.Provider(), id.Types(), id.Names()))
	return builder.String()
}

func formatScope(provider string, types []string, names []string) string {
	if len(types) != len(names) {
		panic(fmt.Sprintf("invalid input: len(%v) != len(%v)", types, names))
	}
	l := len(types)
	segs := make([]string, 1+2*l)
	segs[0] = "/providers/" + provider
	for i := 0; i < l; i++ {
		segs[1+2*i] = types[i]
		segs[1+2*i+1] = names[i]
	}
	return strings.Join(segs, "/")
}
