package storage

import (
	config "teredix/pkg/config"

	mock "github.com/stretchr/testify/mock"

	resource "teredix/pkg/resource"
)

// StorageMock is an autogenerated mock type for the Storage type
type Mock struct {
	mock.Mock
}

// Find provides a mock function with given fields: filter
func (_m *Mock) Find(filter ResourceFilter) ([]resource.Resource, error) {
	ret := _m.Called(filter)

	var r0 []resource.Resource
	if rf, ok := ret.Get(0).(func(ResourceFilter) []resource.Resource); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resource.Resource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ResourceFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Persist provides a mock function with given fields: resources
func (_m *Mock) Persist(resources []resource.Resource) error {
	ret := _m.Called(resources)

	var r0 error
	if rf, ok := ret.Get(0).(func([]resource.Resource) error); ok {
		r0 = rf(resources)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Prepare provides a mock function with given fields:
func (_m *Mock) Prepare() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreRelations provides a mock function with given fields: criteria
func (_m *Mock) StoreRelations(relation config.Relation) error {
	ret := _m.Called(relation)

	var r0 error
	if rf, ok := ret.Get(0).(func(config.Relation) error); ok {
		r0 = rf(relation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRelations provides a mock function with given fields:
func (_m *Mock) GetRelations() ([]map[string]string, error) {
	ret := _m.Called()

	var r0 []map[string]string
	if rf, ok := ret.Get(0).(func() []map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResources provides a mock function with given fields:
func (_m *Mock) GetResources() ([]resource.Resource, error) {
	ret := _m.Called()

	var r0 []resource.Resource
	if rf, ok := ret.Get(0).(func() []resource.Resource); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resource.Resource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
