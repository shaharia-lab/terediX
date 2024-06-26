package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
)

func TestFsScanner_Scan(t *testing.T) {
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
			expectedMetaDataCount: 2,
			expectedMetaDataKeys: []string{
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
			expectedMetaDataCount: 2,
			expectedMetaDataKeys: []string{
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

			sc := config.Source{
				Type:          pkg.SourceTypeFileSystem,
				Configuration: map[string]string{"root_directory": tmpDir},
				Fields:        []string{fileSystemFieldRootDirectory, fileSystemFieldMachineHost},
				Schedule:      "@every 1s",
			}

			mockStorage := new(storage.Mock)
			mockStorage.On("GetNextVersionForResource", "scanner_name", pkg.SourceTypeFileSystem).Return(1, nil)

			fs := FsScanner{}
			_ = fs.Setup("scanner_name", sc, NewScannerDependencies(scheduler.NewGoCron(), mockStorage, &logrus.Logger{}, metrics.NewCollector()))

			RunCommonScannerAssertionTest(t, &fs, tt.expectedResourceCount, tt.expectedMetaDataCount, tt.expectedMetaDataKeys)
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

// Checks if all the keys in the given list exist in the metaData of a Resource
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

// Helper function to check if a single key exists in metaData
func keyExists(resource resource.Resource, key string) bool {
	data := resource.GetMetaData()
	for _, md := range data.Get() {
		if md.Key == key {
			return true
		}
	}
	return false
}
