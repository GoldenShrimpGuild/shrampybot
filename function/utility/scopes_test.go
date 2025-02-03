package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchScope(t *testing.T) {
	testCases := []struct {
		name        string
		userScopes  []string
		targetScope string
		expected    bool
	}{
		{
			name:        "Exact match",
			userScopes:  []string{"admin:users"},
			targetScope: "admin:users",
			expected:    true,
		},
		{
			name:        "Partial match - user scope more specific",
			userScopes:  []string{"admin:users:edit"},
			targetScope: "admin:users",
			expected:    false,
		},
		{
			name:        "Partial match - target scope more specific",
			userScopes:  []string{"admin"},
			targetScope: "admin:users",
			expected:    true,
		},

		{
			name:        "No match",
			userScopes:  []string{"user:view"},
			targetScope: "admin:users",
			expected:    false,
		},
		{
			name:        "Empty user scopes",
			userScopes:  []string{},
			targetScope: "admin:users",
			expected:    false,
		},
		{
			name:        "Empty target scope",
			userScopes:  []string{"admin:users"},
			targetScope: "",
			expected:    false,
		},
		{
			name:        "Multiple user scopes - match",
			userScopes:  []string{"user:view", "admin:users"},
			targetScope: "admin:users",
			expected:    true,
		},
		{
			name:        "Multiple user scopes - no match",
			userScopes:  []string{"user:view", "admin:categories"},
			targetScope: "admin:users",
			expected:    false,
		},
		{
			name:        "Top level scope match",
			userScopes:  []string{"admin"},
			targetScope: "admin",
			expected:    true,
		},
		{
			name:        "Top level scope mismatch",
			userScopes:  []string{"admin"},
			targetScope: "user",
			expected:    false,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MatchScope(tc.userScopes, tc.targetScope)
			assert.Equal(t, tc.expected, result)
		})
	}
}
