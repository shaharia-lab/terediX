package cmd

import "testing"

func Test_validate(t *testing.T) {
	tests := []struct {
		name    string
		cfgFile string
		wantErr bool
	}{
		{
			name:    "valid config file",
			cfgFile: "testdata/valid_config.yaml",
			wantErr: false,
		},
		{
			name:    "invalid config file",
			cfgFile: "testdata/invalid_config.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
