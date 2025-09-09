package namer

import (
	"strings"
	"testing"
)

// Test cases covering various resource naming scenarios
func TestNewResourceName_NoTruncation(t *testing.T) {
	tests := []struct {
		name         string
		baseName     string
		serviceName  string
		resourceType string
		maxLength    int
		expected     string
	}{
		{
			name:         "simple cache instance",
			baseName:     "fullstack",
			serviceName:  "cache",
			resourceType: "instance",
			maxLength:    50,
			expected:     "fullstack-cache-instance",
		},
		{
			name:         "cache firewall",
			baseName:     "fullstack",
			serviceName:  "cache",
			resourceType: "firewall",
			maxLength:    50,
			expected:     "fullstack-cache-firewall",
		},
		{
			name:         "cloudflare edge waf",
			baseName:     "cloudflare",
			serviceName:  "edge",
			resourceType: "waf",
			maxLength:    50,
			expected:     "cloudflare-edge-waf",
		},
		{
			name:         "frontend account",
			baseName:     "fullstack",
			serviceName:  "frontend",
			resourceType: "account",
			maxLength:    50,
			expected:     "fullstack-frontend-account",
		},
		{
			name:         "proxy app bucket",
			baseName:     "my-proxy-app",
			serviceName:  "bucket",
			resourceType: "as-cache",
			maxLength:    50,
			expected:     "my-proxy-app-bucket-as-cache",
		},
		{
			name:         "no resource type",
			baseName:     "fullstack",
			serviceName:  "frontend",
			resourceType: "",
			maxLength:    50,
			expected:     "fullstack-frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer := New(tt.baseName)
			result := namer.NewResourceName(tt.serviceName, tt.resourceType, tt.maxLength)
			if result != tt.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewResourceName_WithTruncation(t *testing.T) {
	tests := []struct {
		name         string
		baseName     string
		serviceName  string
		resourceType string
		maxLength    int
		expected     string
	}{
		{
			name:         "long cloudflare waf ruleset",
			baseName:     "cloudflare-edge-waf",
			serviceName:  "l7-ruleset-ddos",
			resourceType: "managed",
			maxLength:    30,
			expected:     "cloudf-l7-ruleset-ddos-managed",
		},
		{
			name:         "long cache ruleset optimization",
			baseName:     "cloudflare-edge-waf",
			serviceName:  "cache-ruleset",
			resourceType: "optimization",
			maxLength:    25,
			expected:     "cloudflare-cache-r-optimi",
		},
		{
			name:         "long frontend secret accessor",
			baseName:     "fullstack",
			serviceName:  "frontend",
			resourceType: "secret-accessor",
			maxLength:    20,
			expected:     "fulls-fron-secret-a",
		},
		{
			name:         "long zone dns name",
			baseName:     "cloudflare-edge-waf",
			serviceName:  "zone",
			resourceType: "dns",
			maxLength:    15,
			expected:     "cloudf-zone-dns",
		},
		{
			name:         "long latency slo name",
			baseName:     "fullstack",
			serviceName:  "frontend",
			resourceType: "latency-slo",
			maxLength:    18,
			expected:     "fulls-fron-latenc",
		},
		{
			name:         "long proxy firewall rate limit",
			baseName:     "my-proxy-app",
			serviceName:  "firewall",
			resourceType: "rate-limit",
			maxLength:    20,
			expected:     "my-prox-fire-rate-l",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer := New(tt.baseName)
			result := namer.NewResourceName(tt.serviceName, tt.resourceType, tt.maxLength)

			if result != tt.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, tt.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > tt.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), tt.maxLength)
			}
		})
	}
}

func TestNewResourceName_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		baseName     string
		serviceName  string
		resourceType string
		maxLength    int
		expected     string
	}{
		{
			name:         "very short limit with long names",
			baseName:     "fullstack",
			serviceName:  "frontend",
			resourceType: "account",
			maxLength:    8,
			expected:     "fu-fr-ac",
		},
		{
			name:         "minimal limit",
			baseName:     "app",
			serviceName:  "svc",
			resourceType: "res",
			maxLength:    5,
			expected:     "a-s-r",
		},
		{
			name:         "empty resource type with truncation",
			baseName:     "very-long-base-name",
			serviceName:  "service",
			resourceType: "",
			maxLength:    15,
			expected:     "very-lo-service",
		},
		{
			name:         "single character components",
			baseName:     "a",
			serviceName:  "b",
			resourceType: "c",
			maxLength:    10,
			expected:     "a-b-c",
		},
		{
			name:         "hyphenated base name truncation",
			baseName:     "test-fullstack",
			serviceName:  "cache",
			resourceType: "vpc-connector",
			maxLength:    25,
			expected:     "test-cache-vpc-connector",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer := New(tt.baseName)
			result := namer.NewResourceName(tt.serviceName, tt.resourceType, tt.maxLength)

			if result != tt.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, tt.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > tt.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), tt.maxLength)
			}
		})
	}
}

func TestNewResourceName_CommonInfrastructureExamples(t *testing.T) {
	tests := []struct {
		name         string
		baseName     string
		serviceName  string
		resourceType string
		maxLength    int
		expected     string
	}{
		{
			name:         "gcp cloud sql instance",
			baseName:     "production",
			serviceName:  "database",
			resourceType: "instance",
			maxLength:    63, // GCP Cloud SQL limit
			expected:     "production-database-instance",
		},
		{
			name:         "aws s3 bucket",
			baseName:     "company-prod",
			serviceName:  "data-lake",
			resourceType: "bucket",
			maxLength:    63, // AWS S3 bucket limit
			expected:     "company-prod-data-lake-bucket",
		},
		{
			name:         "kubernetes service",
			baseName:     "microservice",
			serviceName:  "user-auth",
			resourceType: "service",
			maxLength:    63, // Kubernetes service name limit
			expected:     "microservice-user-auth-service",
		},
		{
			name:         "azure storage account",
			baseName:     "enterprise",
			serviceName:  "backup",
			resourceType: "storage",
			maxLength:    24, // Azure storage account limit
			expected:     "enterpris-backup-storage",
		},
		{
			name:         "docker container name",
			baseName:     "app",
			serviceName:  "web-server",
			resourceType: "container",
			maxLength:    30,
			expected:     "app-web-server-container",
		},
		{
			name:         "terraform resource",
			baseName:     "infrastructure",
			serviceName:  "networking",
			resourceType: "vpc",
			maxLength:    40,
			expected:     "infrastructure-networking-vpc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer := New(tt.baseName)
			result := namer.NewResourceName(tt.serviceName, tt.resourceType, tt.maxLength)

			if result != tt.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, tt.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > tt.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), tt.maxLength)
			}

			// Verify no leading/trailing hyphens
			if strings.HasPrefix(result, "-") || strings.HasSuffix(result, "-") {
				t.Errorf("NewResourceName() = %s, should not have leading/trailing hyphens", result)
			}

			// Verify no consecutive hyphens
			if strings.Contains(result, "--") {
				t.Errorf("NewResourceName() = %s, should not contain consecutive hyphens", result)
			}
		})
	}
}

// Helper function for minimum calculation
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
