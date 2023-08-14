// Package scanner scans targets
package scanner

import (
	"github.com/shaharia-lab/teredix/pkg/resource"
)

// Scanner interface to build different scanner
type Scanner interface {
	Scan(resourceChannel chan resource.Resource) error
}

func RunScannerForTests(scanner Scanner) []resource.Resource {
	resourceChannel := make(chan resource.Resource)

	var res []resource.Resource

	go func() {
		scanner.Scan(resourceChannel)
		close(resourceChannel)
	}()

	for r := range resourceChannel {
		res = append(res, r)
	}
	return res
}
