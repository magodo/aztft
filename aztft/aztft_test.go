package aztft

import (
	"testing"

	"github.com/magodo/armid"
	"github.com/stretchr/testify/require"
)

func MustParseId(t *testing.T, id string) armid.ResourceId {
	azureId, err := armid.ParseResourceId(id)
	if err != nil {
		t.Fatal(err)
	}
	return azureId
}

func TestQueryType(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectTypes []Type
		expectExact bool
		err         string
	}{
		{
			name:  "invalid id",
			input: "/subscriptions/sub1/resourceGroups/rg1/foos",
			err:   `invalid resource id: scopes should be split by "/providers/"`,
		},
		{
			name:  "resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/resourceGroups/rg1"),
					TFType:  "azurerm_resource_group",
				},
			},
			expectExact: true,
		},
		{
			name:  "resource group (case insensitively)",
			input: "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1"),
					TFType:  "azurerm_resource_group",
				},
			},
			expectExact: true,
		},
		{
			name:  "management group",
			input: "/providers/Microsoft.Management/managementGroups/group1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementGroups/group1"),
					TFType:  "azurerm_management_group",
				},
			},
			expectExact: true,
		},
		{
			name:  "management group (case insensitively)",
			input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1"),
					TFType:  "azurerm_management_group",
				},
			},
			expectExact: true,
		},
		{
			name:  "poliy definition (subscription level)",
			input: "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1"),
					TFType:  "azurerm_policy_definition",
				},
			},
			expectExact: true,
		},
		{
			name:  "policy definitinon (management group level)",
			input: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1"),
					TFType:  "azurerm_policy_definition",
				},
			},
			expectExact: true,
		},
		{
			name:  "policy set definition (subscription level)",
			input: "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1"),
					TFType:  "azurerm_policy_set_definition",
				},
			},
			expectExact: true,
		},
		{
			name:  "policy set definitinon (management group level)",
			input: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1"),
					TFType:  "azurerm_policy_set_definition",
				},
			},
			expectExact: true,
		},
		{
			name:  "backup protection resource",
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"),
					TFType:  "azurerm_backup_protected_file_share",
				},
				{
					AzureId: MustParseId(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"),
					TFType:  "azurerm_backup_protected_vm",
				},
			},
			expectExact: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualType, actualExact, err := QueryType(tt.input, nil)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectTypes, actualType)
			require.Equal(t, tt.expectExact, actualExact)
		})
	}
}

func TestQueryId(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		rt     string
		expect string
		err    string
	}{
		{
			name:  "invalid id",
			input: "/subscriptions/sub1/resourceGroups/rg1/foos",
			rt:    "azurerm_resource_group",
			err:   `parsing id: scopes should be split by "/providers/"`,
		},
		{
			name:   "resource group",
			input:  "/subscriptions/sub1/resourceGroups/rg1",
			rt:     "azurerm_resource_group",
			expect: "/subscriptions/sub1/resourceGroups/rg1",
		},
		{
			name:   "resource group (case insensitively)",
			input:  "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1",
			rt:     "azurerm_resource_group",
			expect: "/subscriptions/SUB1/resourceGroups/RG1",
		},
		{
			name:   "management group",
			input:  "/providers/Microsoft.Management/managementGroups/group1",
			rt:     "azurerm_management_group",
			expect: "/providers/Microsoft.Management/managementGroups/group1",
		},
		{
			name:   "management group (case insensitively)",
			input:  "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1",
			rt:     "azurerm_management_group",
			expect: "/providers/Microsoft.Management/managementGroups/GROUP1",
		},
		{
			name:   "poliy definition (subscription level)",
			input:  "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			rt:     "azurerm_policy_definition",
			expect: "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
		},
		{
			name:   "policy definitinon (management group level)",
			input:  "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			rt:     "azurerm_policy_definition",
			expect: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
		},
		{
			name:   "policy set definition (subscription level)",
			input:  "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			rt:     "azurerm_policy_set_definition",
			expect: "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
		},
		{
			name:   "policy set definitinon (management group level)",
			input:  "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			rt:     "azurerm_policy_set_definition",
			expect: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
		},
		{
			name:   "backup protection resource",
			input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			rt:     "azurerm_backup_protected_vm",
			expect: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := QueryId(tt.input, tt.rt, nil)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
		})
	}
}

func TestQueryTypeAndId(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectTypes []Type
		expectIds   []string
		expectExact bool
		err         string
	}{
		{
			name:  "invalid id",
			input: "/subscriptions/sub1/resourceGroups/rg1/foos",
			err:   `invalid resource id: scopes should be split by "/providers/"`,
		},
		{
			name:  "resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/resourceGroups/rg1"),
					TFType:  "azurerm_resource_group",
				},
			},
			expectIds:   []string{"/subscriptions/sub1/resourceGroups/rg1"},
			expectExact: true,
		},
		{
			name:  "resource group (case insensitively)",
			input: "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1"),
					TFType:  "azurerm_resource_group",
				},
			},
			expectIds:   []string{"/subscriptions/SUB1/resourceGroups/RG1"},
			expectExact: true,
		},
		{
			name:  "management group",
			input: "/providers/Microsoft.Management/managementGroups/group1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementGroups/group1"),
					TFType:  "azurerm_management_group",
				},
			},
			expectIds:   []string{"/providers/Microsoft.Management/managementGroups/group1"},
			expectExact: true,
		},
		{
			name:  "management group (case insensitively)",
			input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1"),
					TFType:  "azurerm_management_group",
				},
			},
			expectIds:   []string{"/providers/Microsoft.Management/managementGroups/GROUP1"},
			expectExact: true,
		},
		{
			name:  "poliy definition (subscription level)",
			input: "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1"),
					TFType:  "azurerm_policy_definition",
				},
			},
			expectIds:   []string{"/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1"},
			expectExact: true,
		},
		{
			name:  "policy definitinon (management group level)",
			input: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1"),
					TFType:  "azurerm_policy_definition",
				},
			},
			expectIds:   []string{"/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1"},
			expectExact: true,
		},
		{
			name:  "policy set definition (subscription level)",
			input: "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1"),
					TFType:  "azurerm_policy_set_definition",
				},
			},
			expectIds:   []string{"/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1"},
			expectExact: true,
		},
		{
			name:  "policy set definitinon (management group level)",
			input: "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1"),
					TFType:  "azurerm_policy_set_definition",
				},
			},
			expectIds:   []string{"/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1"},
			expectExact: true,
		},
		{
			name:  "backup protection resource",
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"),
					TFType:  "azurerm_backup_protected_file_share",
				},
				{
					AzureId: MustParseId(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1"),
					TFType:  "azurerm_backup_protected_vm",
				},
			},
			expectIds: []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			},
			expectExact: false,
		},
		{
			name:  "app service slot virtual network swift connection",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Web/sites/site1/slots/slot1/networkConfig/cfg1",
			expectTypes: []Type{
				{
					AzureId: MustParseId(t, "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Web/sites/site1/slots/slot1/networkConfig/cfg1"),
					TFType:  "azurerm_app_service_slot_virtual_network_swift_connection",
				},
			},
			expectIds:   []string{"/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Web/sites/site1/slots/slot1/config/cfg1"},
			expectExact: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualTypes, actualIds, actualExact, err := QueryTypeAndId(tt.input, nil)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.Equal(t, tt.expectTypes, actualTypes)
			require.Equal(t, tt.expectIds, actualIds)
			require.Equal(t, tt.expectExact, actualExact)
		})
	}
}
