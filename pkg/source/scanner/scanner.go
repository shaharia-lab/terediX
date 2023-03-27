package scanner

import (
	"teredix/pkg/resource"
)

type Scanner interface {
	Scan() []resource.Resource
}
