// Package scanner scans targets
package scanner

import (
	"teredix/pkg/resource"
)

// Scanner interface to build different scanner
type Scanner interface {
	Scan() []resource.Resource
}
