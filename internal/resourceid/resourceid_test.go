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
			input:  TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			input:  SubscriptionId{Id: "sub1"},
			expect: "/subscriptions/sub1",
		},
		{
			name:   "Resource Group",
			input:  ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/sub1/resourceGroups/rg1",
		},
		{
			name:   "Management Group",
			input:  ManagementGroup{Name: "mg1"},
			expect: "/providers/Microsoft.Management/managementGroups/mg1",
		},
		{
			name: "Scoped Resource under tenant",
			input: ScopedResourceId{
				Scope:         TenantId{},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under subscription",
			input: ScopedResourceId{
				Scope:         SubscriptionId{Id: "sub1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under resource group",
			input: ScopedResourceId{
				Scope:         ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under management group",
			input: ScopedResourceId{
				Scope:         ManagementGroup{Name: "mg1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under another scoped resource which under tenant",
			input: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         TenantId{},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under subscription",
			input: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         SubscriptionId{Id: "sub1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under resource group",
			input: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under management group",
			input: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         ManagementGroup{Name: "mg1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
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

func TestResourceId_Parent(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect ResourceId
	}{
		{
			name:   "Tenant",
			input:  TenantId{},
			expect: nil,
		},
		{
			name:   "Subscription",
			input:  SubscriptionId{Id: "sub1"},
			expect: nil,
		},
		{
			name:   "Resource Group",
			input:  ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: nil,
		},
		{
			name:   "Management Group",
			input:  ManagementGroup{Name: "mg1"},
			expect: nil,
		},
		{
			name: "Root Scoped Resource under tenant",
			input: ScopedResourceId{
				Scope:         TenantId{},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos"},
				ResourceNames: []string{"foo1"},
			},
			expect: nil,
		},
		{
			name: "Child Scoped Resource under tenant",
			input: ScopedResourceId{
				Scope:         TenantId{},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
			expect: ScopedResourceId{
				Scope:         TenantId{},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos"},
				ResourceNames: []string{"foo1"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.Parent())
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
			expect: TenantId{},
		},
		{
			name:   "Subscription",
			input:  "/subscriptions/sub1",
			expect: SubscriptionId{Id: "sub1"},
		},
		{
			name:   "Resource Group",
			input:  "/subscriptions/sub1/resourceGroups/rg1",
			expect: ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Case-insensitive for resourceGroups",
			input:  "/subscriptions/sub1/RESOURCEGROUPS/rg1",
			expect: ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Management Group",
			input:  "/providers/Microsoft.Management/managementGroups/mg1",
			expect: ManagementGroup{Name: "mg1"},
		},
		{
			name:   "Case-insensitive for managementGroup",
			input:  "/providers/Microsoft.Management/MANAGEMENTGROUPS/mg1",
			expect: ManagementGroup{Name: "mg1"},
		},
		{
			name:  "Scoped Resource under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: ScopedResourceId{
				Scope:         TenantId{},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: ScopedResourceId{
				Scope:         SubscriptionId{Id: "sub1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: ScopedResourceId{
				Scope:         ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: ScopedResourceId{
				Scope:         ManagementGroup{Name: "mg1"},
				Namespace:     "Microsoft.Foo",
				ResourceTypes: []string{"foos", "bars"},
				ResourceNames: []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         TenantId{},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         SubscriptionId{Id: "sub1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: ScopedResourceId{
				Scope: ScopedResourceId{
					Scope:         ManagementGroup{Name: "mg1"},
					Namespace:     "Microsoft.Foo",
					ResourceTypes: []string{"foos", "bars"},
					ResourceNames: []string{"foo1", "bar1"},
				},
				Namespace:     "Microsoft.Baz",
				ResourceTypes: []string{"bazs"},
				ResourceNames: []string{"baz1"},
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
