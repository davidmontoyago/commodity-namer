package namer_test

import (
	"strings"
	"testing"

	namer "github.com/davidmontoyago/commodity-namer"
)

// Test cases covering various resource naming scenarios
func TestNewResourceName_NoTruncation(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			n := namer.New(testCase.baseName)
			result := n.NewResourceName(testCase.serviceName, testCase.resourceType, testCase.maxLength)
			if result != testCase.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, testCase.expected)
			}
		})
	}
}

func TestNewResourceName_WithTruncation(t *testing.T) {
	t.Parallel()

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
		{
			name:         "long backend processor service account",
			baseName:     "my-prod-stack",
			serviceName:  "backend-processor",
			resourceType: "service-account",
			maxLength:    30,
			expected:     "my-prod-backend-pr-service-a",
		},
		{
			name:         "long backend processor service account",
			baseName:     "my-prod-stack",
			serviceName:  "ingestor",
			resourceType: "generic-service",
			maxLength:    25,
			expected:     "my-prod-inges-generic-s",
		},
		{
			name:         "long require https",
			baseName:     "my-prod-stack",
			serviceName:  "require-https",
			resourceType: "",
			maxLength:    20,
			expected:     "my-pro-require-https",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			n := namer.New(testCase.baseName)
			result := n.NewResourceName(testCase.serviceName, testCase.resourceType, testCase.maxLength)

			if result != testCase.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, testCase.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > testCase.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), testCase.maxLength)
			}
		})
	}
}

func TestNewResourceName_EdgeCases(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			n := namer.New(testCase.baseName)
			result := n.NewResourceName(testCase.serviceName, testCase.resourceType, testCase.maxLength)

			if result != testCase.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, testCase.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > testCase.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), testCase.maxLength)
			}
		})
	}
}

func TestNewResourceName_CommonInfrastructureExamples(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			n := namer.New(testCase.baseName)
			result := n.NewResourceName(testCase.serviceName, testCase.resourceType, testCase.maxLength)

			if result != testCase.expected {
				t.Errorf("NewResourceName() = %v, want %v", result, testCase.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > testCase.maxLength {
				t.Errorf("NewResourceName() length = %d, want <= %d",
					len(result), testCase.maxLength)
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

func TestNewResourceName_InvalidName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		baseName     string
		serviceName  string
		resourceType string
		maxLength    int
		expectPanic  bool
		description  string
	}{
		{
			name:         "name starts with hyphen",
			baseName:     "-invalid",
			serviceName:  "service",
			resourceType: "type",
			maxLength:    50,
			expectPanic:  true,
			description:  "should panic when name starts with hyphen",
		},
		{
			name:         "name ends with hyphen",
			baseName:     "invalid",
			serviceName:  "service",
			resourceType: "type-",
			maxLength:    50,
			expectPanic:  true,
			description:  "should panic when name ends with hyphen",
		},
		{
			name:         "name starts with digit",
			baseName:     "9invalid",
			serviceName:  "service",
			resourceType: "type",
			maxLength:    50,
			expectPanic:  true,
			description:  "should panic when name starts with digit",
		},
		{
			name:         "name with uppercase letters",
			baseName:     "Invalid",
			serviceName:  "service",
			resourceType: "type",
			maxLength:    50,
			expectPanic:  true,
			description:  "should panic when name contains uppercase letters",
		},
		{
			name:         "name with special characters",
			baseName:     "invalid_name",
			serviceName:  "service",
			resourceType: "type",
			maxLength:    50,
			expectPanic:  true,
			description:  "should panic when name contains special characters",
		},
		{
			name:         "name exceeds length limit",
			baseName:     "a-very-very-very-very-very-very-very-very-long-base-name",
			serviceName:  "service",
			resourceType: "type",
			maxLength:    70,
			expectPanic:  true,
			description:  "should panic when name exceeds RFC 1035 limit of 63 characters",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); r != nil {
					if !testCase.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
					// Expected panic occurred
				} else {
					if testCase.expectPanic {
						t.Errorf("Expected panic but none occurred for: %s", testCase.description)
					}
				}
			}()

			n := namer.New(testCase.baseName)
			result := n.NewResourceName(testCase.serviceName, testCase.resourceType, testCase.maxLength)

			// If we reach here and expectPanic is true, the test should fail
			if testCase.expectPanic {
				t.Errorf("Expected panic but got result: %s", result)
			}
		})
	}
}
