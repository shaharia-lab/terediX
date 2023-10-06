package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMissingKeys(t *testing.T) {
	tests := []struct {
		name        string
		metadata    []MetaData
		keysToCheck []string
		expected    []string
	}{
		{
			name: "All keys exist",
			metadata: []MetaData{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
			keysToCheck: []string{"key1", "key2", "key3"},
			expected:    []string{},
		},
		{
			name: "Some keys missing",
			metadata: []MetaData{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			keysToCheck: []string{"key1", "key3"},
			expected:    []string{"key3"},
		},
		{
			name:        "No keys exist in metadata",
			metadata:    []MetaData{},
			keysToCheck: []string{"key1", "key2"},
			expected:    []string{"key1", "key2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MetaDataLists{
				data: tt.metadata,
			}
			got := ml.FindMissingKeys(tt.keysToCheck)

			assert.Equal(t, len(got), len(tt.expected))
		})
	}
}
