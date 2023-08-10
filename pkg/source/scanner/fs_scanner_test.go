package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shahariaazam/teredix/pkg"
	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/shahariaazam/teredix/pkg/resource"

	"github.com/stretchr/testify/assert"
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
				"filex.txt": "file1 content",
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
					Name:       fmt.Sprintf("%s/%s", tmpDir, "filex.txt"),
					ExternalID: fmt.Sprintf("%s/%s", tmpDir, "filex.txt"),
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
				"filex.txt":        "file1 content",
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
					Name:       fmt.Sprintf("%s/%s", tmpDir, "filex.txt"),
					ExternalID: fmt.Sprintf("%s/%s", tmpDir, "filex.txt"),
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

			resourceChannel := make(chan resource.Resource, len(tt.expected))
			var res []resource.Resource

			go func() {
				// Create an FsScanner for the temporary directory and scan it
				scanner := NewFsScanner("scanner_name", tmpDir, []string{"rootDirectory", "machineHost"})
				scanner.Scan(resourceChannel)
				close(resourceChannel)
			}()

			for r := range resourceChannel {
				res = append(res, r)
			}

			for k, r := range res {
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

func TestFsScanner_ScanV2(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name                  string
		fsSource              config.Source
		files                 map[string]string
		expectedResourceCount int
		expectedMetaDataCount int
		expectedMetaDataKeys  []string
	}{
		{
			name: "single file",
			fsSource: config.Source{
				Type: pkg.SourceTypeFileSystem,
				Configuration: map[string]string{
					"root_directory": tmpDir,
				},
				Fields: []string{fileSystemFieldRootDirectory, fileSystemFieldMachineHost},
			},
			files: map[string]string{
				"filex.txt": "file1 content",
			},
			expectedResourceCount: 2,
			expectedMetaDataCount: 6,
			expectedMetaDataKeys: []string{
				"Machine-Host",
				"Root-Directory",
				"Scanner-Label",
				fileSystemFieldMachineHost,
				fileSystemFieldRootDirectory,
			},
		},
		{
			name: "nested directory",
			fsSource: config.Source{
				Type: pkg.SourceTypeFileSystem,
				Configuration: map[string]string{
					"root_directory": tmpDir,
				},
				Fields: []string{fileSystemFieldRootDirectory, fileSystemFieldMachineHost},
			},
			files: map[string]string{
				"dir1/nested1.txt": "content",
				"filex.txt":        "file1 content",
			},
			expectedResourceCount: 3,
			expectedMetaDataCount: 6,
			expectedMetaDataKeys: []string{
				"Machine-Host",
				"Root-Directory",
				"Scanner-Label",
				fileSystemFieldMachineHost,
				fileSystemFieldRootDirectory,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create files
			err := generateTmpTestFiles(tmpDir, tt.files)
			if err != nil {
				t.Errorf(err.Error())
			}

			resourceChannel := make(chan resource.Resource)
			var res []resource.Resource

			go func() {
				// Create an FsScanner for the temporary directory and scan it
				scanner := NewFsScanner("scanner_name", tmpDir, []string{"rootDirectory", "machineHost"})
				scanner.Scan(resourceChannel)
				close(resourceChannel)
			}()

			for r := range resourceChannel {
				res = append(res, r)
			}

			assert.Equal(t, tt.expectedResourceCount, len(res), fmt.Sprintf("expected %d resource, but got %d resource", tt.expectedResourceCount, len(res)))
			assert.Equal(t, tt.expectedMetaDataCount, len(res[0].MetaData))

			for k, v := range res {
				exists, missingKeys := checkKeysInMetaData(v, tt.expectedMetaDataKeys)
				if !exists {
					t.Errorf("Metadata missing. Missing keys [%d]: %v", k, missingKeys)
				}
			}
		})
	}
}

func generateTmpTestFiles(targetDirectory string, files map[string]string) error {
	for filename, content := range files {
		// Split the path into directory and file name
		dirPath, _ := filepath.Split(filename)
		// Create the directory hierarchy
		if err := os.MkdirAll(filepath.Join(targetDirectory, dirPath), 0755); err != nil {
			return fmt.Errorf("Error creating directory hierarchy: %v", err)
		}
		// Write the file
		filePath := filepath.Join(targetDirectory, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("Error creating test file: %v", err)
		}
	}
	return nil
}

// Checks if all the keys in the given list exist in the MetaData of a Resource
// Returns a boolean indicating if all keys exist and a slice of missing keys
func checkKeysInMetaData(resource resource.Resource, keys []string) (bool, []string) {
	missingKeys := []string{}

	for _, key := range keys {
		if !keyExists(resource, key) {
			missingKeys = append(missingKeys, key)
		}
	}

	return len(missingKeys) == 0, missingKeys
}

// Helper function to check if a single key exists in MetaData
func keyExists(resource resource.Resource, key string) bool {
	for _, md := range resource.MetaData {
		if md.Key == key {
			return true
		}
	}
	return false
}
