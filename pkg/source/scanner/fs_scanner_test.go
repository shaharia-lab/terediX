package scanner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"teredix/pkg/resource"
	"testing"
)

func TestFsScanner_Scan(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		files          map[string]string
		attachMetaData map[string]string
		expected       []resource.Resource
	}{
		{
			name: "single file",
			files: map[string]string{
				"file1.txt": "file1 content",
			},
			attachMetaData: map[string]string{"key1": "value1", "key2": "value2"},
			expected: []resource.Resource{
				{
					Kind:       "FileDirectory",
					UUID:       "xxxx",
					Name:       tmpDir,
					ExternalID: tmpDir,
					MetaData: []resource.MetaData{
						{
							Key:   "Scanner",
							Value: "scanner_name",
						},
						{
							Key:   "key1",
							Value: "value1",
						},
						{
							Key:   "key2",
							Value: "value2",
						},
					},
					RelatedWith: []resource.Resource{},
				},
				{
					Kind:       "FilePath",
					UUID:       "xxxx",
					Name:       fmt.Sprintf("%s/%s", tmpDir, "file1.txt"),
					ExternalID: fmt.Sprintf("%s/%s", tmpDir, "file1.txt"),
					RelatedWith: []resource.Resource{
						{
							Kind: "FileDirectory",
							UUID: "xxxx",
							Name: tmpDir,
						},
					},
					MetaData: []resource.MetaData{
						{
							Key:   "Scanner",
							Value: "scanner_name",
						},
						{
							Key:   "key1",
							Value: "value1",
						},
						{
							Key:   "key2",
							Value: "value2",
						},
					},
				},
			},
		},
		{
			name: "nested directories",
			files: map[string]string{
				"dir1/nested1.txt": "content",
				"file1.txt":        "file1 content",
			},
			attachMetaData: map[string]string{},
			expected: []resource.Resource{
				{
					Kind:        "FileDirectory",
					UUID:        "xxxx",
					Name:        tmpDir,
					ExternalID:  tmpDir,
					RelatedWith: []resource.Resource{},
					MetaData: []resource.MetaData{
						{
							Key:   "Scanner",
							Value: "scanner_name",
						},
					},
				},
				{
					Kind:       "FilePath",
					UUID:       "xxxx",
					Name:       fmt.Sprintf("%s/%s", tmpDir, "dir1/nested1.txt"),
					ExternalID: fmt.Sprintf("%s/%s", tmpDir, "dir1/nested1.txt"),
					RelatedWith: []resource.Resource{
						{
							Kind: "FileDirectory",
							UUID: "xxxx",
							Name: tmpDir,
						},
					},
					MetaData: []resource.MetaData{
						{
							Key:   "Scanner",
							Value: "scanner_name",
						},
					},
				},
				{
					Kind:       "FilePath",
					UUID:       "xxxx",
					Name:       fmt.Sprintf("%s/%s", tmpDir, "file1.txt"),
					ExternalID: fmt.Sprintf("%s/%s", tmpDir, "file1.txt"),
					RelatedWith: []resource.Resource{
						{
							Kind: "FileDirectory",
							UUID: "xxxx",
							Name: tmpDir,
						},
					},
					MetaData: []resource.MetaData{
						{
							Key:   "Scanner",
							Value: "scanner_name",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create files
			for filename, content := range tt.files {
				// Split the path into directory and file name
				dirPath, _ := filepath.Split(filename)
				// Create the directory hierarchy
				if err := os.MkdirAll(filepath.Join(tmpDir, dirPath), 0755); err != nil {
					t.Fatalf("Error creating directory hierarchy: %v", err)
				}
				// Write the file
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Error creating test file: %v", err)
				}
			}

			// Create an FsScanner for the temporary directory and scan it
			scanner := NewFsScanner("scanner_name", tmpDir, tt.attachMetaData)
			resources := scanner.Scan()
			assert.Equal(t, len(tt.expected), len(resources))
			for k, r := range resources {
				assert.Equal(t, tt.expected[k].Kind, r.Kind)
				assert.Equal(t, tt.expected[k].Name, r.Name)
				assert.Equal(t, tt.expected[k].ExternalID, r.ExternalID)
				assert.Equal(t, "scanner_name", r.FindMetaValue("Scanner"))
				assert.Equal(t, len(tt.expected[k].RelatedWith), len(r.RelatedWith))
				if len(tt.expected[k].RelatedWith) > 0 {
					assert.Equal(t, tt.expected[k].RelatedWith[0].Name, r.RelatedWith[0].Name)
				}
			}
		})
	}
}
