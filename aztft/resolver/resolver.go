package resolver

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type Resolver struct {
	cred azcore.TokenCredential
}

func (r *Resolver) Resolve(routeScopeStr string, parentScopeStr string) (string, error) {
	switch routeScopeStr {
	}
	return "", fmt.Errorf("Unknown route scope: %s", routeScopeStr)
}
