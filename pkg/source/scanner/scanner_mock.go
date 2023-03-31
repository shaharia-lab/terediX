// Package scanner scans targets
package scanner

import (
	"teredix/pkg/resource"

	"github.com/stretchr/testify/mock"
)

// Mock is an autogenerated mock type for the ScannerMock type
type Mock struct {
	mock.Mock
}

// Scan provides a mock function with given fields: resourceChannel
func (_m *Mock) Scan(resourceChannel chan resource.Resource) error {
	ret := _m.Called(resourceChannel)

	var r0 error
	if rf, ok := ret.Get(0).(func(chan resource.Resource) error); ok {
		r0 = rf(resourceChannel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
