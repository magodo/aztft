package aztft

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryType(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect []string
		err    string
	}{
		{
			name:  "invalid id",
			input: "/subscriptions/sub1/resourceGroups/rg1/foos",
			err:   `invalid resource id: scopes should be split by "/providers/"`,
		},
		{
			name:   "resource group",
			input:  "/subscriptions/sub1/resourceGroups/rg1",
			expect: []string{"azurerm_resource_group"},
		},
		{
			name:   "resource group (case insensitively)",
			input:  "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1",
			expect: []string{"azurerm_resource_group"},
		},
		{
			name:   "management group",
			input:  "/providers/Microsoft.Management/managementGroups/group1",
			expect: []string{"azurerm_management_group"},
		},
		{
			name:   "management group (case insensitively)",
			input:  "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1",
			expect: []string{"azurerm_management_group"},
		},
		{
			name:   "poliy definition (subscription level)",
			input:  "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expect: []string{"azurerm_policy_definition"},
		},
		{
			name:   "policy definitinon (management group level)",
			input:  "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expect: []string{"azurerm_policy_definition"},
		},
		{
			name:   "policy set definition (subscription level)",
			input:  "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expect: []string{"azurerm_policy_set_definition"},
		},
		{
			name:   "policy set definitinon (management group level)",
			input:  "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expect: []string{"azurerm_policy_set_definition"},
		},
		{
			name:   "backup protection resource",
			input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			expect: []string{"azurerm_backup_protected_file_share", "azurerm_backup_protected_vm"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := QueryType(tt.input, false)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.Equal(t, tt.expect, actual)
		})
	}
}

func TestQueryTypeAndId(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectRts   []string
		expectSpecs []string
		err         string
	}{
		{
			name:  "invalid id",
			input: "/subscriptions/sub1/resourceGroups/rg1/foos",
			err:   `invalid resource id: scopes should be split by "/providers/"`,
		},
		{
			name:        "resource group",
			input:       "/subscriptions/sub1/resourceGroups/rg1",
			expectRts:   []string{"azurerm_resource_group"},
			expectSpecs: []string{"/subscriptions/sub1/resourceGroups/rg1"},
		},
		{
			name:        "resource group (case insensitively)",
			input:       "/SUBSCRIPTIONS/SUB1/RESOURCEGROUPS/RG1",
			expectRts:   []string{"azurerm_resource_group"},
			expectSpecs: []string{"/subscriptions/SUB1/resourceGroups/RG1"},
		},
		{
			name:        "management group",
			input:       "/providers/Microsoft.Management/managementGroups/group1",
			expectRts:   []string{"azurerm_management_group"},
			expectSpecs: []string{"/providers/Microsoft.Management/managementGroups/group1"},
		},
		{
			name:        "management group (case insensitively)",
			input:       "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1",
			expectRts:   []string{"azurerm_management_group"},
			expectSpecs: []string{"/providers/Microsoft.Management/managementGroups/GROUP1"},
		},
		{
			name:        "poliy definition (subscription level)",
			input:       "/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectRts:   []string{"azurerm_policy_definition"},
			expectSpecs: []string{"/subscriptions/sub1/providers/Microsoft.Authorization/policyDefinitions/policy1"},
		},
		{
			name:        "policy definitinon (management group level)",
			input:       "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1",
			expectRts:   []string{"azurerm_policy_definition"},
			expectSpecs: []string{"/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policyDefinitions/policy1"},
		},
		{
			name:        "policy set definition (subscription level)",
			input:       "/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectRts:   []string{"azurerm_policy_set_definition"},
			expectSpecs: []string{"/subscriptions/sub1/providers/Microsoft.Authorization/policySetDefinitions/policy1"},
		},
		{
			name:        "policy set definitinon (management group level)",
			input:       "/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1",
			expectRts:   []string{"azurerm_policy_set_definition"},
			expectSpecs: []string{"/providers/Microsoft.Management/managementgroups/grp1/providers/Microsoft.Authorization/policySetDefinitions/policy1"},
		},
		{
			name:      "backup protection resource",
			input:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			expectRts: []string{"azurerm_backup_protected_file_share", "azurerm_backup_protected_vm"},
			expectSpecs: []string{
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;group1;vm1/protectedItems/vm;iaasvmcontainerv2;group1;vm1",
			},
		},
		{
			name:        "app service slot virtual network swift connection",
			input:       "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Web/sites/site1/slots/slot1/networkConfig/cfg1",
			expectRts:   []string{"azurerm_app_service_slot_virtual_network_swift_connection"},
			expectSpecs: []string{"/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Web/sites/site1/slots/slot1/config/cfg1"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualRts, actualSpecs, err := QueryTypeAndId(tt.input, false)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.Equal(t, tt.expectRts, actualRts)
			require.Equal(t, tt.expectSpecs, actualSpecs)
		})
	}
}
