package resolve

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type springApmsResolver struct{}

func (springApmsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_spring_cloud_dynatrace_application_performance_monitoring",
		"azurerm_spring_cloud_application_insights_application_performance_monitoring",
	}
}

func (springApmsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	return "", fmt.Errorf("can't resolve for Spring Apms as the client is not supporting it yet")
}
