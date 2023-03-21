package scanner

import (
	"infrastructure-discovery/pkg/resource"
)

type Scanner interface {
	Scan() []resource.Resource
}
