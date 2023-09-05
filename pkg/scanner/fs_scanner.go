// Package scanner scans targets
package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shaharia-lab/teredix/pkg/resource"
)

const (
	fileSystemFieldRootDirectory = "rootDirectory"
	fileSystemFieldMachineHost   = "machineHost"
)

// FsScanner store configuration for file system scanner
type FsScanner struct {
	name          string
	rootDirectory string
	fields        []string
}

// File represent file information
type File struct {
	Path string
}

// NewFsScanner construct new file system scanner
func NewFsScanner(name, rootDirectory string, fields []string) *FsScanner {
	return &FsScanner{name: name, rootDirectory: rootDirectory, fields: fields}
}

// Scan scans the file system
func (s *FsScanner) Scan(resourceChannel chan resource.Resource) error {
	files, err := s.listFilesRecursive(s.rootDirectory)
	if err != nil {
		return nil
	}

	//rootResource := resource.NewResourceV1("FileDirectory", util.GenerateUUID(), s.rootDirectory, s.rootDirectory, s.name)
	rootResource := resource.NewResource("FileDirectory", s.name, s.rootDirectory, s.name, 1)

	mappings := map[string]func() string{
		fileSystemFieldRootDirectory: func() string { return s.rootDirectory },
		fileSystemFieldMachineHost: func() string {
			hostname, err := os.Hostname()
			if err != nil {
				hostname = ""
			}
			return hostname
		},
	}

	resourceMeta := NewFieldMapper(mappings, nil, s.fields).getResourceMetaData()
	rootResource.AddMetaData(resourceMeta)

	resourceChannel <- rootResource

	for _, f := range files {
		//nr := resource.NewResourceV1("FilePath", util.GenerateUUID(), f.Path, f.Path, s.name)
		nr := resource.NewResource("FilePath", s.name, f.Path, s.name, 1)
		nr.AddRelation(rootResource)
		nr.AddMetaData(resourceMeta)

		resourceChannel <- nr
	}

	return nil
}

func (s *FsScanner) listFilesRecursive(path string) ([]File, error) {
	var files []File
	info, err := os.Stat(path)
	if err != nil {
		return files, err
	}

	if !info.IsDir() {
		return append(files, File{Path: path}), nil
	}

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		fileOrDirPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			subFiles, err := s.listFilesRecursive(fileOrDirPath)
			if err != nil {
				return files, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, File{Path: fileOrDirPath})
		}
	}
	return files, nil
}
