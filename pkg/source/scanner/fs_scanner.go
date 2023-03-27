package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"teredix/pkg/resource"
	"teredix/pkg/util"
)

type FsScanner struct {
	name          string
	rootDirectory string
	metaData      map[string]string
}

type File struct {
	Path string
}

func NewFsScanner(name, rootDirectory string, metaData map[string]string) FsScanner {
	return FsScanner{name: name, rootDirectory: rootDirectory, metaData: metaData}
}

func (s *FsScanner) Scan() []resource.Resource {
	files, err := s.listFilesRecursive(s.rootDirectory)
	if err != nil {
		return nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	var r []resource.Resource

	rootResource := resource.NewResource("FileDirectory", util.GenerateUUID(), s.rootDirectory, s.rootDirectory, s.name)
	for k, v := range s.metaData {
		rootResource.AddMetaData(k, v)
	}

	rootResource.AddMetaData("Machine-Host", hostname)
	rootResource.AddMetaData("Root-Directory", s.rootDirectory)
	rootResource.AddMetaData("Scanner-Label", s.name)

	r = append(r, rootResource)

	for _, f := range files {
		nr := resource.NewResource("FilePath", util.GenerateUUID(), f.Path, f.Path, s.name)
		nr.AddRelation(rootResource)
		for k, v := range s.metaData {
			nr.AddMetaData(k, v)
		}

		nr.AddMetaData("Machine-Host", hostname)
		nr.AddMetaData("Root-Directory", s.rootDirectory)
		nr.AddMetaData("Scanner-Label", s.name)

		r = append(r, nr)
	}

	return r
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
