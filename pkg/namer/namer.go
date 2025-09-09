// Package namer allows generating consistent resource names with length limits.
package namer

import (
	"fmt"
	"math"
)

// Namer provides consistent resource naming with length constraints
type Namer struct {
	baseName string
}

// New creates a new Namer instance with the given base name
func New(baseName string) *Namer {
	return &Namer{baseName: baseName}
}

// NewResourceName generates a consistent resource name with length limits.
func (e *Namer) NewResourceName(resourceName, resourceType string, maxLength int) string {
	var name string
	if resourceType == "" {
		name = fmt.Sprintf("%s-%s", e.baseName, resourceName)
	} else {
		name = fmt.Sprintf("%s-%s-%s", e.baseName, resourceName, resourceType)
	}

	if len(name) <= maxLength {
		return name
	}

	surplus := len(name) - maxLength
	name = e.truncateResourceName(resourceName, resourceType, surplus, maxLength)

	return name
}

// truncateResourceName handles the complex logic for truncating resource names.
func (e *Namer) truncateResourceName(serviceName, resourceType string, surplus, maxLength int) string {
	mainComponentLength := len(e.baseName)
	if mainComponentLength > surplus {
		return e.truncateMainComponent(serviceName, resourceType, surplus)
	}

	return e.proportionalTruncate(serviceName, resourceType, maxLength)
}

// truncateMainComponent truncates the main component name when it's long enough.
func (e *Namer) truncateMainComponent(serviceName, resourceType string, surplus int) string {
	truncatedMainComponent := e.baseName[:len(e.baseName)-surplus]
	truncatedMainComponent = trimTrailingHyphen(truncatedMainComponent)

	if resourceType == "" {
		return fmt.Sprintf("%s-%s", truncatedMainComponent, serviceName)
	}

	return fmt.Sprintf("%s-%s-%s", truncatedMainComponent, serviceName, resourceType)
}

// proportionalTruncate applies proportional truncation when main component is too short.
func (e *Namer) proportionalTruncate(serviceName, resourceType string, maxLength int) string {
	originalLength := len(fmt.Sprintf("%s-%s-%s", e.baseName, serviceName, resourceType))
	if resourceType == "" {
		originalLength = len(fmt.Sprintf("%s-%s", e.baseName, serviceName))
	}

	truncateFactorFloat := float64(maxLength) / float64(originalLength)
	truncateFactor := math.Floor(truncateFactorFloat*100) / 100

	mainComponentLength := int(math.Floor(float64(len(e.baseName)) * truncateFactor))
	serviceNameLength := int(math.Floor(float64(len(serviceName)) * truncateFactor))
	resourceTypeLength := int(math.Floor(float64(len(resourceType)) * truncateFactor))

	// Truncate each component and remove trailing hyphens
	truncatedBaseName := trimTrailingHyphen(e.baseName[:mainComponentLength])
	truncatedServiceName := trimTrailingHyphen(serviceName[:serviceNameLength])
	truncatedResourceType := trimTrailingHyphen(resourceType[:resourceTypeLength])

	return join(truncatedBaseName, truncatedServiceName, truncatedResourceType)
}

func join(truncatedBaseName string, truncatedServiceName string, truncatedResourceType string) string {
	if truncatedResourceType == "" {
		return fmt.Sprintf("%s-%s", truncatedBaseName, truncatedServiceName)
	}

	return fmt.Sprintf("%s-%s-%s", truncatedBaseName, truncatedServiceName, truncatedResourceType)
}

// trimTrailingHyphen removes trailing hyphens from a string component
func trimTrailingHyphen(component string) string {
	for len(component) > 0 && component[len(component)-1] == '-' {
		component = component[:len(component)-1]
	}

	return component
}
