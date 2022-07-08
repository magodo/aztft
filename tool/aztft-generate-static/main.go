package main

/// This program generate the mapping from the Azure document.
/// The generated mapping only covers the management plane, if a resource has a data plane ID in the meanwhile, that mapping info have to be manually added.
/// For Terraform resource ids in the Azure document that failed to parse (e.g. due to it is a data plane id), or is a synthetic (e.g. association resource) id, an empty mapping item is generated.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/magodo/aztft/internal/resmap"

	"errors"

	"github.com/magodo/armid"
)

var (
	ErrMalformedImportSpec error = errors.New("malformed import spec")
	ErrDataPlaneId         error = errors.New("data plane id")
	ErrSyntheticId         error = errors.New("synthetic id")
	ErrParseIdFailed       error = errors.New("failed to parse id")

	knownParseErrors = []error{
		ErrMalformedImportSpec,
		ErrDataPlaneId,
		ErrSyntheticId,
		ErrParseIdFailed,
	}

	ErrDuplicateImportSpec error = errors.New("duplicate import spec")
)

type HardCodedTypeInfo struct {
	mapItem   *resmap.TF2ARMIdMapItem
	caughtErr error
	caught    bool
}

var HardcodedTypes = map[string]*HardCodedTypeInfo{
	// Associations are skipped as they are beyond the goal of this tool
	"azurerm_app_service_certificate_binding":                                        {caughtErr: ErrSyntheticId},
	"azurerm_management_group_subscription_association":                              {caughtErr: ErrParseIdFailed},
	"azurerm_network_interface_security_group_association":                           {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_application_gateway_backend_address_pool_association": {caughtErr: ErrSyntheticId},
	"azurerm_virtual_desktop_workspace_application_group_association":                {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_application_security_group_association":               {caughtErr: ErrSyntheticId},
	"azurerm_nat_gateway_public_ip_association":                                      {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_nat_rule_association":                                 {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_backend_address_pool_association":                     {caughtErr: ErrSyntheticId},
	"azurerm_nat_gateway_public_ip_prefix_association":                               {caughtErr: ErrSyntheticId},
	"azurerm_disk_pool_managed_disk_attachment":                                      {caughtErr: ErrSyntheticId},

	// Data plane only resources
	"azurerm_storage_share_file":                                     {caughtErr: ErrDataPlaneId},
	"azurerm_storage_share_directory":                                {caughtErr: ErrDataPlaneId},
	"azurerm_storage_table_entity":                                   {caughtErr: ErrDataPlaneId},
	"azurerm_storage_blob":                                           {caughtErr: ErrDataPlaneId},
	"azurerm_key_vault_certificate":                                  {caughtErr: ErrDataPlaneId},
	"azurerm_storage_data_lake_gen2_path":                            {caughtErr: ErrDataPlaneId},
	"azurerm_storage_data_lake_gen2_filesystem":                      {caughtErr: ErrDataPlaneId},
	"azurerm_key_vault_managed_storage_account":                      {caughtErr: ErrDataPlaneId},
	"azurerm_key_vault_managed_storage_account_sas_token_definition": {caughtErr: ErrDataPlaneId},
	"azurerm_key_vault_certificate_issuer":                           {caughtErr: ErrDataPlaneId},
	"azurerm_synapse_role_assignment":                                {caughtErr: ErrSyntheticId},

	// This is not a azure resource, but an operation like abstract resource. Skip it.
	"azurerm_resource_provider_registration": {caughtErr: ErrParseIdFailed},
	// This represents a property in disk pool iscsi target resource, which means it doesn't have a resource ID in Azure. Just ignore it.
	"azurerm_disk_pool_iscsi_target_lun": {caughtErr: ErrSyntheticId},

	"azurerm_key_vault_access_policy": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "objectId", "applicationId"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.KeyVault/vaults/objectId/applicationId"},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_security_center_server_vulnerability_assessment": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups/Microsoft.Compute/virtualMachines"},
				Provider:     "Microsoft.Security",
				Types:        []string{"serverVulnerabilityAssessments"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.Compute/virtualMachines/Microsoft.Security/serverVulnerabilityAssessments"},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_backup_protected_vm": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.RecoveryServices",
				Types:        []string{"vaults", "backupFabrics", "protectionContainers", "protectedItems"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_backup_protected_file_share": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.RecoveryServices",
				Types:        []string{"vaults", "backupFabrics", "protectionContainers", "protectedItems"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_backup_container_storage_account": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.RecoveryServices",
				Types:        []string{"vaults", "backupFabrics", "protectionContainers"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_policy_definition": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions", "/Microsoft.Management/managementGroups"},
				Provider:     "Microsoft.Authorization",
				Types:        []string{"policyDefinitions"},
				ImportSpecs: []string{
					"/subscriptions/Microsoft.Authorization/policyDefinitions",
					"/Microsoft.Management/managementgroups/Microsoft.Authorization/policyDefinitions",
				},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_policy_set_definition": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions", "/Microsoft.Management/managementGroups"},
				Provider:     "Microsoft.Authorization",
				Types:        []string{"policySetDefinitions"},
				ImportSpecs: []string{
					"/subscriptions/Microsoft.Authorization/policySetDefinitions",
					"/Microsoft.Management/managementgroups/Microsoft.Authorization/policySetDefinitions",
				},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_app_configuration_feature": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.AppConfiguration",
				Types:        []string{"configurationStores", "AppConfigurationFeature", "Label"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.AppConfiguration/configurationStores/AppConfigurationFeature/Label"},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_app_configuration_key": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.AppConfiguration",
				Types:        []string{"configurationStores", "AppConfigurationKey", "Label"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.AppConfiguration/configurationStores/AppConfigurationKey/Label"},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_monitor_diagnostic_setting": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Insights",
				Types:        []string{"diagnosticSettings"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_storage_object_replication": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "objectReplicationPolicies"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_resource_policy_assignment": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Authorization",
				Types:        []string{"policyAssignments"},
			},
		},
		caughtErr: ErrParseIdFailed,
	},
	"azurerm_role_definition": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Authorization",
				Types:        []string{"roleDefinitions"},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_storage_share": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "fileServices", "shares"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_container": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "blobServices", "containers"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_queue": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "queueServices", "queues"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_table": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "tableServices", "tables"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_key": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "keys"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_secret": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "secrets"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
}

const usage = `aztft-generate-static <provider root dir>`

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	rootDir := os.Args[1]
	rDir := path.Join(rootDir, "website", "docs", "r")
	dir, err := os.Open(rDir)
	if err != nil {
		log.Fatalf("failed to read directory %s: %v\n", rDir, err)
	}
	entries, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatalf("failed to read directory entries under %s: %v\n", rDir, err)
	}
	dir.Close()

	m := map[string]armid.ResourceId{}

	for _, entry := range entries {
		p := path.Join(rDir, entry)
		f, err := os.Open(p)
		if err != nil {
			log.Fatalf("failed to open %s: %v\n", p, err)
		}
		scanner := bufio.NewScanner(f)
	ScanFileLoop:
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "terraform import") && !strings.HasPrefix(line, "$ terraform import") {
				continue
			}
			line = line[strings.Index(line, "terraform import"):]
			rtype, id, err := parse(line)
			if err != nil {
				if HardcodedTypes[rtype] == nil {
					log.Printf("%s new parse error: %v\n", rtype, err)
					continue
				}
				// Skip the error if it is already caught
				for _, kerr := range knownParseErrors {
					if errors.Is(err, kerr) {
						if HardcodedTypes[rtype].caughtErr == kerr {
							HardcodedTypes[rtype].caught = true
							continue ScanFileLoop
						}
					}
				}

				log.Fatalf("%s parse error: %v\n", rtype, err)
			}

			if _, ok := m[rtype]; ok {
				if HardcodedTypes[rtype] == nil {
					log.Printf("%s new duplicate import spec found\n", rtype)
					continue
				}
				// Skip this if the duplication is already caught
				if HardcodedTypes[rtype].caughtErr == ErrDuplicateImportSpec {
					HardcodedTypes[rtype].caught = true
					continue
				}

				log.Fatalf("%s duplicate import spec found\n", rtype)
			}
			m[rtype] = id
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("reading %s: %v", p, err)
		}
	}

	// Ensure all the caught errors are really caught
	for rtype, err := range HardcodedTypes {
		if !err.caught {
			log.Fatalf("Expect catch error for %s, but didn't", rtype)
		}
	}

	mapItems := resmap.TF2ARMIdMapType{}
	for rtype, id := range m {
		var scopes []string
		if _, ok := id.(armid.RootScope); !ok {
			scopes = []string{id.ParentScope().ScopeString()}
		}
		mapItems[rtype] = resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: scopes,
				Provider:     id.Provider(),
				Types:        id.Types(),
				ImportSpecs:  []string{id.ScopeString()},
			},
		}
	}
	for rtype, item := range HardcodedTypes {
		if item.mapItem != nil {
			mapItems[rtype] = *item.mapItem
		}
	}

	b, err := json.MarshalIndent(mapItems, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func parse(line string) (string, armid.ResourceId, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return "", nil, fmt.Errorf("%s: %w", line, ErrMalformedImportSpec)
	}
	addr, idRaw := fields[2], fields[3]
	rtype, _, ok := strings.Cut(addr, ".")
	if !ok {
		return "", nil, fmt.Errorf("%s: malformed resource address", addr)
	}

	if v, err := strconv.Unquote(idRaw); err == nil {
		idRaw = v
	}

	if strings.HasPrefix(idRaw, "https://") {
		return rtype, nil, ErrDataPlaneId
	}

	// Return an empty TF2ARMIdMapItem for the synthetic resources, which are mostly binding/association resources.
	if strings.ContainsAny(idRaw, ";|") {
		return rtype, nil, ErrSyntheticId
	}

	id, err := armid.ParseResourceId(idRaw)
	if err != nil {
		return rtype, nil, fmt.Errorf("%w: %v", ErrParseIdFailed, err)
	}

	return rtype, id, nil
}
