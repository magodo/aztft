package main

/// This program generate the mapping from the Azure document.
/// The generated mapping only covers the management plane, if a resource has a data plane ID in the meanwhile, that mapping info have to be manually added.
/// For Terraform resource ids in the Azure document that failed to parse (e.g. due to it is a data plane id), or is a synthetic (e.g. association resource) id, an empty mapping item is generated.

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"errors"

	"github.com/magodo/aztft/internal/resourceid"
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

type ErrorInfo struct {
	Err    error
	Caught bool
}

var CaughtErrors = map[string]*ErrorInfo{
	"azurerm_app_service_certificate_binding":                                        {Err: ErrSyntheticId},
	"azurerm_management_group_subscription_association":                              {Err: ErrParseIdFailed},
	"azurerm_network_interface_security_group_association":                           {Err: ErrSyntheticId},
	"azurerm_backup_protected_vm":                                                    {Err: ErrSyntheticId},
	"azurerm_policy_set_definition":                                                  {Err: ErrDuplicateImportSpec},
	"azurerm_storage_share":                                                          {Err: ErrDataPlaneId},
	"azurerm_network_interface_application_gateway_backend_address_pool_association": {Err: ErrSyntheticId},
	"azurerm_key_vault_certificate":                                                  {Err: ErrDataPlaneId},
	// "azurerm_eventgrid_event_subscription":                                           {Err: ErrMalformedImportSpec},
	"azurerm_app_configuration_feature":                               {Err: ErrDuplicateImportSpec},
	"azurerm_policy_definition":                                       {Err: ErrDuplicateImportSpec},
	"azurerm_monitor_diagnostic_setting":                              {Err: ErrSyntheticId},
	"azurerm_virtual_desktop_workspace_application_group_association": {Err: ErrSyntheticId},
	"azurerm_storage_data_lake_gen2_path":                             {Err: ErrDataPlaneId},
	// "azurerm_app_service_source_control_token":                                       ErrParseIdFailed,
	"azurerm_storage_share_file":                                       {Err: ErrDataPlaneId},
	"azurerm_storage_container":                                        {Err: ErrDataPlaneId},
	"azurerm_synapse_role_assignment":                                  {Err: ErrSyntheticId},
	"azurerm_disk_pool_iscsi_target_lun":                               {Err: ErrSyntheticId},
	"azurerm_network_interface_application_security_group_association": {Err: ErrSyntheticId},
	"azurerm_policy_remediation":                                       {Err: ErrDuplicateImportSpec},
	"azurerm_key_vault_key":                                            {Err: ErrDataPlaneId},
	// "azurerm_orchestrated_virtual_machine_scale_set":                   ErrParseIdFailed,
	"azurerm_backup_protected_file_share":                            {Err: ErrSyntheticId},
	"azurerm_storage_queue":                                          {Err: ErrDataPlaneId},
	"azurerm_storage_blob":                                           {Err: ErrDataPlaneId},
	"azurerm_nat_gateway_public_ip_association":                      {Err: ErrSyntheticId},
	"azurerm_app_configuration_key":                                  {Err: ErrDuplicateImportSpec},
	"azurerm_storage_object_replication":                             {Err: ErrSyntheticId},
	"azurerm_key_vault_secret":                                       {Err: ErrDataPlaneId},
	"azurerm_linux_web_app":                                          {Err: ErrDuplicateImportSpec},
	"azurerm_security_center_server_vulnerability_assessment":        {Err: ErrDuplicateImportSpec},
	"azurerm_storage_share_directory":                                {Err: ErrDataPlaneId},
	"azurerm_key_vault_managed_storage_account":                      {Err: ErrDataPlaneId},
	"azurerm_resource_policy_assignment":                             {Err: ErrParseIdFailed},
	"azurerm_storage_data_lake_gen2_filesystem":                      {Err: ErrDataPlaneId},
	"azurerm_storage_table":                                          {Err: ErrDataPlaneId},
	"azurerm_backup_container_storage_account":                       {Err: ErrSyntheticId},
	"azurerm_key_vault_access_policy":                                {Err: ErrDuplicateImportSpec},
	"azurerm_disk_pool_managed_disk_attachment":                      {Err: ErrSyntheticId},
	"azurerm_network_interface_nat_rule_association":                 {Err: ErrSyntheticId},
	"azurerm_resource_provider_registration":                         {Err: ErrParseIdFailed},
	"azurerm_hpc_cache_blob_target":                                  {Err: ErrDuplicateImportSpec},
	"azurerm_storage_table_entity":                                   {Err: ErrDataPlaneId},
	"azurerm_network_interface_backend_address_pool_association":     {Err: ErrSyntheticId},
	"azurerm_nat_gateway_public_ip_prefix_association":               {Err: ErrSyntheticId},
	"azurerm_key_vault_managed_storage_account_sas_token_definition": {Err: ErrDataPlaneId},
	"azurerm_key_vault_certificate_issuer":                           {Err: ErrDataPlaneId},
	"azurerm_role_definition":                                        {Err: ErrSyntheticId},
}

const usage = `aztft-generate-static <provider root dir>`

type MapItem struct {
	ManagementPlane MapManagementPlane `json:"management_plane,omitempty"`
	DataPlane       MapDataPlane       `json:"data_plane,omitempty"`
}

type MapManagementPlane struct {
	Scope ScopeManagementPlane `json:"scope"`
	Type  string               `json:"type"`
}

type MapDataPlane struct {
	Scope ScopeDataPlane `json:"scope"`
	Type  string         `json:"type"`
}

type ScopeManagementPlane string

const (
	ManagementPlaneScopeRoot            ScopeManagementPlane = "root"
	ManagementPlaneScopeTenant          ScopeManagementPlane = "tenant"
	ManagementPlaneScopeManagementGroup ScopeManagementPlane = "management_group"
	ManagementPlaneScopeSubscription    ScopeManagementPlane = "subscription"
	ManagementPlaneScopeResourceGroup   ScopeManagementPlane = "resource_group"
)

type ScopeDataPlane string

// const (
// 	DataPlaneScopeKeyVault            ScopeDataPlane = "keyvault"
// 	DataPlaneScopeStorageAccountBlob  ScopeDataPlane = "storage_account_blob"
// 	DataPlaneScopeStorageAccountTable ScopeDataPlane = "storage_account_table"
// 	DataPlaneScopeStorageAccountFile  ScopeDataPlane = "storage_account_file"
// 	DataPlaneScopeStorageAccountQueue ScopeDataPlane = "storage_account_queue"
// )

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

	m := map[string]resourceid.ResourceId{}

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
			if !strings.HasPrefix(line, "terraform import") {
				continue
			}
			rtype, id, err := parse(line)
			if err != nil {
				// Skip the error if it is already caught
				for _, kerr := range knownParseErrors {
					if errors.Is(err, kerr) && CaughtErrors[rtype].Err == kerr {
						CaughtErrors[rtype].Caught = true
						continue ScanFileLoop
					}
				}

				log.Fatalf("%s parse error: %v\n", rtype, err)
				// log.Printf("%s parse error: %v\n", rtype, err)
				// continue
			}
			if _, ok := m[rtype]; ok {
				// Skip this if the duplication is already caught
				if CaughtErrors[rtype].Err == ErrDuplicateImportSpec {
					CaughtErrors[rtype].Caught = true
					continue
				}

				log.Fatalf("%s duplicate import spec found\n", rtype)
				// log.Printf("%s duplicate import spec found\n", rtype)
				// continue
			}
			m[rtype] = id
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("reading %s: %v", p, err)
		}
	}

	// Ensure all the caught errors are really caught
	for rtype, err := range CaughtErrors {
		if !err.Caught {
			log.Fatalf("Expect catch error for %s, but didn't", rtype)
		}
	}
}

func parse(line string) (string, resourceid.ResourceId, error) {
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

	// Return an empty MapItem for the synthetic resources, which are mostly binding/association resources.
	if strings.ContainsAny(idRaw, ";|") {
		return rtype, nil, ErrSyntheticId
	}

	id, err := resourceid.ParseResourceId(idRaw)
	if err != nil {
		return rtype, nil, fmt.Errorf("%w: %v", ErrParseIdFailed, err)
	}

	return rtype, id, nil
	// // Identify the scope
	// var scope ScopeManagementPlane
	// if id.ParentScope() == nil {
	// 	scope = ManagementPlaneScopeRoot
	// } else {
	// 	switch id.ParentScope().(type) {
	// 	case resourceid.TenantId:
	// 		scope = ManagementPlaneScopeTenant
	// 	case resourceid.ManagementGroup:
	// 		scope = ManagementPlaneScopeManagementGroup
	// 	case resourceid.SubscriptionId:
	// 		scope = ManagementPlaneScopeSubscription
	// 	case resourceid.ResourceGroup:
	// 		scope = ManagementPlaneScopeResourceGroup
	// 	}
	// }
}
