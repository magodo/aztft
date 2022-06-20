package resourceid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceId_String(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions/sub1",
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/sub1/resourceGroups/rg1",
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: "/providers/Microsoft.Management/managementGroups/mg1",
		},
		{
			name: "Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under subscription",
			input: &ScopedResourceId{
				AttrParentScope: &SubscriptionId{Id: "sub1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under management group",
			input: &ScopedResourceId{
				AttrParentScope: &ManagementGroup{Name: "mg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under another scoped resource which under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under subscription",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &SubscriptionId{Id: "sub1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under management group",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ManagementGroup{Name: "mg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.String())
		})
	}
}

func TestResourceId_RootScope(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect RootScope
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: &TenantId{},
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: &SubscriptionId{Id: "sub1"},
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name: "Root Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: &TenantId{},
		},
		{
			name: "Child Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &TenantId{},
		},
		{
			name: "Child Scoped Resource under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.RootScope())
		})
	}
}

func TestResourceId_Parent(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect ResourceId
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: nil,
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: nil,
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: nil,
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: nil,
		},
		{
			name: "Root Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: nil,
		},
		{
			name: "Child Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.Parent())
		})
	}
}

func TestResourceId_Equal(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		oid    ResourceId
		expect bool
	}{
		{
			name:   "Tenant equals to Tenant",
			id:     &TenantId{},
			oid:    &TenantId{},
			expect: true,
		},
		{
			name:   "Tenant not equals to Subscription",
			id:     &TenantId{},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: false,
		},
		{
			name:   "Subscription equals to subscription with same id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: true,
		},
		{
			name:   "Subscription not equals to subscription with different id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub2"},
			expect: false,
		},
		{
			name:   "Resource Group equals to Resource Group with same subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: true,
		},
		{
			name:   "Resource Group not equals to Resource Group with different subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub2", Name: "rg2"},
			expect: false,
		},
		{
			name:   "Management Group equals to Management Group with same name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg1"},
			expect: true,
		},
		{
			name:   "Management Group not equals to Management Group with different name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg2"},
			expect: false,
		},
		{
			name: "Root Scoped Resource under tenant equals to itself",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: true,
		},
		{
			name: "Root Scoped Resource under tenant not equals to different resource id",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"bars"},
				AttrNames:       []string{"bar1"},
			},
			expect: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.Equal(tt.oid))
		})
	}
}

func TestResourceId_ScopeEqual(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		oid    ResourceId
		expect bool
	}{
		{
			name:   "Tenant equals scope to Tenant",
			id:     &TenantId{},
			oid:    &TenantId{},
			expect: true,
		},
		{
			name:   "Tenant not equals scope to Subscription",
			id:     &TenantId{},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: false,
		},
		{
			name:   "Subscription equals scope to subscription with different id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub2"},
			expect: true,
		},
		{
			name:   "Resource Group equals scope to Resource Group with different subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub2", Name: "rg2"},
			expect: true,
		},
		{
			name:   "Management Group equals scope to Management Group with different name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg2"},
			expect: true,
		},
		{
			name: "Root Scoped Resource under tenant equals scopes to different sub-type name",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo2"},
			},
			expect: true,
		},
		{
			name: "Parent Scoped Resource under tenant not equals scopes to different sub-type type",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"bars"},
				AttrNames:       []string{"bar1"},
			},
			expect: false,
		},
		{
			name: "Parent Scoped Resource under tenant not equals scopes to its child Scoped Resource",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.ScopeEqual(tt.oid))
		})
	}
}

func TestResourceId_ScopeString(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			id:     &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			id:     &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions",
		},
		{
			name:   "Resource Group",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/resourceGroups",
		},
		{
			name:   "Management Group",
			id:     &ManagementGroup{Name: "mg1"},
			expect: "/Microsoft.Management/managementGroups",
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Child Scoped Resource under resource group with",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos", "bars"},
				AttrNames:    []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/resourceGroups/Microsoft.Foo/foos/bars",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			expect: "/subscriptions/resourceGroups/Microsoft.Foo/foos",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.ScopeString())
		})
	}
}

func TestResourceId_RouteScopeString(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			id:     &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			id:     &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions",
		},
		{
			name:   "Resource Group",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/resourceGroups",
		},
		{
			name:   "Management Group",
			id:     &ManagementGroup{Name: "mg1"},
			expect: "/Microsoft.Management/managementGroups",
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Child Scoped Resource under resource group with",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos", "bars"},
				AttrNames:    []string{"foo1", "bar1"},
			},
			expect: "/Microsoft.Foo/foos/bars",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.RouteScopeString())
		})
	}
}

func TestScopedResourceId_Normalize(t *testing.T) {
	cases := []struct {
		name     string
		id       ScopedResourceId
		scopeStr string
		expect   string
		err      string
	}{
		{
			name: "Mismatch scope string",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/Microsoft.Bar/foos",
			err:      `mismatch scope string ("/Microsoft.Bar/foos") for id "/providers/Microsoft.Foo/foos/foo1"`,
		},
		{
			name: "Root Scoped Resource under tenant",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.Foo",
				AttrTypes:       []string{"FOOS"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/providers/microsoft.foo/foos/foo1",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "MICROSOFT.Foo",
				AttrTypes:    []string{"FOOS"},
				AttrNames:    []string{"foo1"},
			},
			scopeStr: "/subscriptions/resourceGroups/microsoft.foo/foos",
			expect:   "/subscriptions/sub1/resourceGroups/rg1/providers/microsoft.foo/foos/foo1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			err := id.Normalize(tt.scopeStr)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id.String())
		})
	}
}

func TestScopedResourceId_NormalizeRouteScope(t *testing.T) {
	cases := []struct {
		name     string
		id       ScopedResourceId
		scopeStr string
		expect   string
		err      string
	}{
		{
			name: "Mismatch scope string",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/Microsoft.Bar/foos",
			err:      `mismatch route scope string ("/Microsoft.Bar/foos") for id "/providers/Microsoft.Foo/foos/foo1"`,
		},
		{
			name: "Root Scoped Resource under tenant",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.Foo",
				AttrTypes:       []string{"FOOS"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/providers/microsoft.foo/foos/foo1",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "MICROSOFT.Foo",
				AttrTypes:    []string{"FOOS"},
				AttrNames:    []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/subscriptions/sub1/resourceGroups/rg1/providers/microsoft.foo/foos/foo1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			err := id.NormalizeRouteScope(tt.scopeStr)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id.String())
		})
	}
}

func TestParseResourceId(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect ResourceId
		err    string
	}{
		{
			name:   "Tenant",
			input:  "/",
			expect: &TenantId{},
		},
		{
			name:   "Subscription",
			input:  "/subscriptions/sub1",
			expect: &SubscriptionId{Id: "sub1"},
		},
		{
			name:   "Resource Group",
			input:  "/subscriptions/sub1/resourceGroups/rg1",
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Case-insensitive for resourceGroups",
			input:  "/SUBSCRIPTIONS/sub1/RESOURCEGROUPS/rg1",
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Management Group",
			input:  "/providers/Microsoft.Management/managementGroups/mg1",
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name:   "Case-insensitive for managementGroup",
			input:  "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/mg1",
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name:  "Scoped Resource under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Case-insensitiev Scoped Resource under tenant",
			input: "/PROVIDERS/MICROSOFT.FOO/FOOS/foo1/BARS/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.FOO",
				AttrTypes:       []string{"FOOS", "BARS"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &SubscriptionId{Id: "sub1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &ManagementGroup{Name: "mg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &SubscriptionId{Id: "sub1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ManagementGroup{Name: "mg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "empty string",
			input: "",
			err:   `id should start with "/"`,
		},
		{
			name:  "id not starts with /",
			input: "foo",
			err:   `id should start with "/"`,
		},
		{
			name:  `id ends with "/"`,
			input: "/providers/",
			err:   `empty segment found behind 2th "/"`,
		},
		{
			name:  `id has empty segment in the middle "/"`,
			input: "/providers/Microsoft.Foo/foos//foo1",
			err:   `empty segment found behind 4th "/"`,
		},
		{
			name:  "invalid scope behind tenant scope",
			input: "/foo",
			err:   `scopes should be split by "/providers/"`,
		},
		{
			name:  "invalid scope behind subscription scope",
			input: "/subscriptions/sub1/foo",
			err:   `scopes should be split by "/providers/"`,
		},
		{
			name:  "invalid scope behind resource group scope",
			input: "/subscriptions/sub1/resourceGroups/rg1/foo",
			err:   `scopes should be split by "/providers/"`,
		},
		{
			name:  "invalid scope behind management group scope",
			input: "/providers.Management/managementGroups/mg1/foo",
			err:   `scopes should be split by "/providers/"`,
		},
		{
			name:  `missing provider namespace segment`,
			input: "/providers",
			err:   `missing provider namespace segment`,
		},
		{
			name:  `missing sub-type type`,
			input: "/providers/Microsoft.Foo",
			err:   `missing sub-type type`,
		},
		{
			name:  `missing sub-type name`,
			input: "/providers/Microsoft.Foo/foos",
			err:   `missing sub-type name`,
		},
		{
			name:  `missing sub-type name in child`,
			input: "/providers/Microsoft.Foo/foos/foo1/bars",
			err:   `missing sub-type name`,
		},
		{
			name:  `no sub-type in a scope`,
			input: "/providers/Microsoft.Foo/providers/Microsoft.Bar",
			err:   `missing sub-type type`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseResourceId(tt.input)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id)
		})
	}
}
