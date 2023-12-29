package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *cobra.Command
	}{
		{
			name:    "Version is empty",
			version: "",
			want: &cobra.Command{
				Version: "",
				Short:   "Discover and Explore your tech resources",
				Long:    "Discover and Explore your tech resources",
			},
		},
		{
			name:    "Version is not empty",
			version: "1.0.0",
			want: &cobra.Command{
				Version: "1.0.0",
				Short:   "Discover and Explore your tech resources",
				Long:    "Discover and Explore your tech resources",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRootCmd(tt.version)
			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.Short, got.Short)
			assert.Equal(t, tt.want.Long, got.Long)

			// Check if there are exactly 3 commands added
			assert.Equal(t, 3, len(got.Commands()))

			expectedCommands := []*cobra.Command{
				NewDiscoverCommand(),
				NewRelationCommand(),
				NewValidateCommand(),
			}
			for i, cmd := range got.Commands() {
				assert.Equal(t, expectedCommands[i].Use, cmd.Use)
			}

		})
	}
}
