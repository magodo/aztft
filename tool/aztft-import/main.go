package main

/// This program generate the mapping from the Azure document.

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
	// Unknown
	// (not supported)
	"azurerm_pim_active_role_assignment":   {caughtErr: ErrSyntheticId},
	"azurerm_pim_eligible_role_assignment": {caughtErr: ErrSyntheticId},

	// Property-like resources
	// (not supported)
	"azurerm_app_service_certificate_binding":                         {caughtErr: ErrSyntheticId},
	"azurerm_app_service_source_control_token":                        {caughtErr: ErrParseIdFailed},
	"azurerm_private_endpoint_application_security_group_association": {caughtErr: ErrSyntheticId},
	"azurerm_virtual_machine_gallery_application_assignment":          {caughtErr: ErrSyntheticId},
	"azurerm_virtual_desktop_scaling_plan_host_pool_association":      {caughtErr: ErrSyntheticId},
	"azurerm_communication_service_email_domain_association":          {caughtErr: ErrSyntheticId},
	//"azurerm_management_group_subscription_association": {}, // Just not supported

	// Data plane resources
	// (not supported)
	"azurerm_key_vault_managed_hardware_security_module_role_definition": {
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_managed_hardware_security_module_role_assignment": {
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_managed_hardware_security_module_key": {
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_managed_hardware_security_module_key_rotation_policy": {
		caughtErr: ErrDataPlaneId,
	},

	// (supported)
	"azurerm_network_interface_security_group_association":                           {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_application_gateway_backend_address_pool_association": {caughtErr: ErrSyntheticId},
	"azurerm_virtual_desktop_workspace_application_group_association":                {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_application_security_group_association":               {caughtErr: ErrSyntheticId},
	"azurerm_nat_gateway_public_ip_association":                                      {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_nat_rule_association":                                 {caughtErr: ErrSyntheticId},
	"azurerm_network_interface_backend_address_pool_association":                     {caughtErr: ErrSyntheticId},
	"azurerm_nat_gateway_public_ip_prefix_association":                               {caughtErr: ErrSyntheticId},
	"azurerm_chaos_studio_target": {
		caughtErr: ErrParseIdFailed,
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Chaos",
				Types:        []string{"targets"},
			},
		},
	},
	"azurerm_chaos_studio_capability": {
		caughtErr: ErrParseIdFailed,
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Chaos",
				Types:        []string{"targets", "capabilities"},
			},
		},
	},
	"azurerm_api_management_api": {
		caughtErr: ErrSyntheticId,
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.ApiManagement",
				Types:        []string{"service", "apis"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.ApiManagement/service/apis"},
			},
		},
	},

	// Data plane only resources, we use pesudo resource id patterns
	"azurerm_key_vault_certificate": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "certificates"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_certificate_issuer": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "certificates", "issuers"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_certificate_contacts": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "certificates", "contacts"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_managed_storage_account": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "storage"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_key_vault_managed_storage_account_sas_token_definition": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "storage", "sas"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_blob": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "blobServices", "containers", "blobs"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_share_file": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "fileServices", "shares", "files"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_share_directory": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "fileServices", "shares", "directories"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_table_entity": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "tableServices", "tables", "partitionKeys", "rowKeys"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_data_lake_gen2_path": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "dfs", "paths"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_storage_data_lake_gen2_filesystem": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Storage",
				Types:        []string{"storageAccounts", "dfs"},
			},
		},
		caughtErr: ErrDataPlaneId,
	},
	"azurerm_synapse_role_assignment": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Synapse",
				Types:        []string{"workspaces", "roleAssignments"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.Synapse/workspaces"},
			},
		},
		caughtErr: ErrSyntheticId,
	},

	// Normal resources, but encounter issues during import
	"azurerm_key_vault_access_policy": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.KeyVault",
				Types:        []string{"vaults", "objectId"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.KeyVault/vaults/objectId"},
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
	"azurerm_app_configuration_feature": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.AppConfiguration",
				Types:        []string{"configurationStores", "AppConfigurationFeature", "Label"},
				ImportSpecs:  []string{"/subscriptions/resourceGroups/Microsoft.AppConfiguration/configurationStores/AppConfigurationFeature/Label"},
			},
		},
		caughtErr: ErrDataPlaneId,
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
		caughtErr: ErrDataPlaneId,
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
	"azurerm_network_manager_deployment": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Network",
				Types: []string{
					"networkManagers",
					"locations",
					"types",
				},
				ImportSpecs: []string{
					"/subscriptions/resourceGroups/Microsoft.Network/networkManagers",
				},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_role_management_policy": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{resmap.ScopeAny},
				Provider:     "Microsoft.Authorization",
				Types: []string{
					"roleManagementPolicies",
				},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_automation_job_schedule": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.Automation",
				Types: []string{
					"automationAccounts",
					"jobSchedules",
				},
			},
		},
		caughtErr: ErrSyntheticId,
	},
	"azurerm_log_analytics_linked_service": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.OperationalInsights",
				Types: []string{
					"workspaces",
					"linkedServices",
				},
			},
		},
		caughtErr: ErrDuplicateImportSpec,
	},
	"azurerm_postgresql_flexible_server_virtual_endpoint": {
		mapItem: &resmap.TF2ARMIdMapItem{
			ManagementPlane: &resmap.MapManagementPlane{
				ParentScopes: []string{"/subscriptions/resourceGroups"},
				Provider:     "Microsoft.DBforPostgreSQL",
				Types: []string{
					"flexibleServers",
					"virtualEndpoints",
				},
				ImportSpecs: []string{
					"/subscriptions/resourceGroups/Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints",
				},
			},
		},
		caughtErr: ErrSyntheticId,
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
