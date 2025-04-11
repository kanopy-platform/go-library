package okta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFilterStringWithGroupNames(t *testing.T) {

	tests := []struct {
		name   string
		input  []string
		output string
	}{
		{
			name:   "single group name",
			input:  []string{"group1"},
			output: "profile.name eq \"group1\"",
		},
		{
			name:   "multiple group names",
			input:  []string{"group1", "group2"},
			output: "profile.name eq \"group1\" or profile.name eq \"group2\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := toFilterString(test.input)
			assert.Equal(t, test.output, result, "Expected filter string does not match the result")
		})
	}
}

func TestBuildFilterNameBatches(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		batchSize int
		expected  []string
	}{
		{
			name:      "single batch",
			input:     []string{"group1", "group2", "group3"},
			batchSize: 3,
			expected:  []string{"profile.name eq \"group1\" or profile.name eq \"group2\" or profile.name eq \"group3\""},
		},
		{
			name:      "multiple batches",
			input:     []string{"group1", "group2", "group3", "group4", "group5"},
			batchSize: 2,
			expected:  []string{"profile.name eq \"group1\" or profile.name eq \"group2\"", "profile.name eq \"group3\" or profile.name eq \"group4\"", "profile.name eq \"group5\""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := buildFilterNameBatches(test.input, test.batchSize)
			assert.Equal(t, test.expected, result, "Expected batches do not match the result")
		})
	}
}
