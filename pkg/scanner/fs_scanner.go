// Package scanner scans targets
package scanner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
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
	schedule      string
	storage       storage.Storage
	logger        *logrus.Logger
	metrics       *metrics.Collector
}

// File represent file information
type File struct {
	Path string
}

// Setup setup file system scanner
func (s *FsScanner) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	s.name = name
	s.rootDirectory = cfg.Configuration["root_directory"]
	s.fields = cfg.Fields
	s.schedule = cfg.Schedule
	s.storage = dependencies.GetStorage()
	s.logger = dependencies.GetLogger()
	s.metrics = dependencies.GetMetrics()

	s.logger.WithFields(logrus.Fields{
		"scanner_name": s.name,
		"scanner_kind": s.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

// GetName return name
func (s *FsScanner) GetName() string {
	return s.name
}

// GetSchedule return schedule
func (s *FsScanner) GetSchedule() string {
	return s.schedule
}

// GetKind return resource kind
func (s *FsScanner) GetKind() string {
	return pkg.ResourceKindFileSystem
}

// Scan scans the file system
func (s *FsScanner) Scan(resourceChannel chan resource.Resource) error {
	s.metrics.CollectTotalScannerJobStatusCount(s.name, s.GetKind(), "running")
	nextVersion, err := s.storage.GetNextVersionForResource(s.name, pkg.ResourceKindFileSystem)
	if err != nil {
		s.metrics.CollectTotalScannerJobStatusCount(s.name, s.GetKind(), "failed")
		s.logger.WithFields(logrus.Fields{
			"scanner_name": s.name,
			"scanner_kind": s.GetKind(),
		}).WithError(err).Error("Unable to get next version for resource")

		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	files, err := s.listFilesRecursive(s.rootDirectory)
	if err != nil {
		s.metrics.CollectTotalScannerJobStatusCount(s.name, s.GetKind(), "failed")
		return nil
	}

	//rootResource := resource.NewResourceV1("FileDirectory", util.GenerateUUID(), s.rootDirectory, s.rootDirectory, s.name)
	rootResource := resource.NewResource(pkg.ResourceKindFileSystem, s.name, s.rootDirectory, s.name, nextVersion)

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

	totalResourceDiscovered := 1

	resourceMeta := NewFieldMapper(mappings, nil, s.fields).getResourceMetaData()
	rootResource.AddMetaData(resourceMeta)

	resourceChannel <- rootResource
	for _, f := range files {
		//nr := resource.NewResourceV1("FilePath", util.GenerateUUID(), f.Path, f.Path, s.name)
		nr := resource.NewResource(pkg.ResourceKindFileSystem, s.name, f.Path, s.name, nextVersion)
		nr.AddRelation(rootResource)
		nr.AddMetaData(resourceMeta)

		resourceChannel <- nr

		totalResourceDiscovered++
	}

	s.logger.WithFields(logrus.Fields{
		"scanner_name":              s.name,
		"scanner_kind":              s.GetKind(),
		"total_resource_discovered": totalResourceDiscovered,
	}).Info("scan completed")

	s.metrics.CollectTotalScannerJobStatusCount(s.name, s.GetKind(), "finished")
	s.metrics.CollectTotalResourceDiscoveredByScanner(s.name, s.GetKind(), strconv.Itoa(nextVersion), float64(totalResourceDiscovered))
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
