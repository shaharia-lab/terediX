// Package scanner scans targets
package scanner

import (
	"github.com/shaharia-lab/teredix/pkg/resource"
)

// Scanner interface to build different scanner
type Scanner interface {
	Scan(resourceChannel chan resource.Resource) error
}
