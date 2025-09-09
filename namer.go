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
func New(baseName string) Namer {
	return Namer{baseName: baseName}
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

// truncateResourceName truncates and handles max length constraints.
func (e *Namer) truncateResourceName(resourceName, resourceType string, surplus, maxLength int) string {
	mainComponentLength := len(e.baseName)
	if mainComponentLength > surplus {
		return e.truncateMainComponent(resourceName, resourceType, surplus)
	}

	return e.proportionalTruncate(resourceName, resourceType, maxLength)
}

// truncateMainComponent truncates the main component name when it's long enough.
func (e *Namer) truncateMainComponent(resourceName, resourceType string, surplus int) string {
	truncatedMainComponent := e.baseName[:len(e.baseName)-surplus]
	truncatedMainComponent = trimTrailingHyphen(truncatedMainComponent)

	if resourceType == "" {
		return fmt.Sprintf("%s-%s", truncatedMainComponent, resourceName)
	}

	return fmt.Sprintf("%s-%s-%s", truncatedMainComponent, resourceName, resourceType)
}

// proportionalTruncate applies proportional truncation when main component is too short.
func (e *Namer) proportionalTruncate(resourceName, resourceType string, maxLength int) string {
	originalLength := len(join(e.baseName, resourceName, resourceType))

	truncateFactorFloat := float64(maxLength) / float64(originalLength)
	truncateFactor := math.Floor(truncateFactorFloat*100) / 100

	mainComponentLength := int(math.Floor(float64(len(e.baseName)) * truncateFactor))
	resourceNameLength := int(math.Floor(float64(len(resourceName)) * truncateFactor))
	resourceTypeLength := int(math.Floor(float64(len(resourceType)) * truncateFactor))

	// Truncate each component and remove trailing hyphens
	truncatedBaseName := trimTrailingHyphen(e.baseName[:mainComponentLength])
	truncatedResourceName := trimTrailingHyphen(resourceName[:resourceNameLength])
	truncatedResourceType := trimTrailingHyphen(resourceType[:resourceTypeLength])

	return join(truncatedBaseName, truncatedResourceName, truncatedResourceType)
}

// join composes the base name, resource name and resource type in the final format.
func join(baseName string, resourceName string, resourceType string) string {
	if resourceType == "" {
		return fmt.Sprintf("%s-%s", baseName, resourceName)
	}

	return fmt.Sprintf("%s-%s-%s", baseName, resourceName, resourceType)
}

// trimTrailingHyphen removes trailing hyphens from a string component
func trimTrailingHyphen(component string) string {
	for len(component) > 0 && component[len(component)-1] == '-' {
		component = component[:len(component)-1]
	}

	return component
}
