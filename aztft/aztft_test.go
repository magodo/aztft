package aztft

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
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
			actual, err := Query(tt.input, false)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.Equal(t, tt.expect, actual)
		})
	}
}
