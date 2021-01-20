// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	dynamic "k8s.io/client-go/dynamic"

	genericclioptions "k8s.io/cli-runtime/pkg/genericclioptions"

	kubernetes "k8s.io/client-go/kubernetes"

	mock "github.com/stretchr/testify/mock"

	resource "k8s.io/cli-runtime/pkg/resource"

	rest "k8s.io/client-go/rest"
)

// Factory is an autogenerated mock type for the Factory type
type Factory struct {
	mock.Mock
}

// Builder provides a mock function with given fields:
func (_m *Factory) Builder() *resource.Builder {
	ret := _m.Called()

	var r0 *resource.Builder
	if rf, ok := ret.Get(0).(func() *resource.Builder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resource.Builder)
		}
	}

	return r0
}

// DynamicClient provides a mock function with given fields:
func (_m *Factory) DynamicClient() (dynamic.Interface, error) {
	ret := _m.Called()

	var r0 dynamic.Interface
	if rf, ok := ret.Get(0).(func() dynamic.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dynamic.Interface)
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

// GetOverrideFlags provides a mock function with given fields:
func (_m *Factory) GetOverrideFlags() (*genericclioptions.ConfigFlags, error) {
	ret := _m.Called()

	var r0 *genericclioptions.ConfigFlags
	if rf, ok := ret.Get(0).(func() *genericclioptions.ConfigFlags); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*genericclioptions.ConfigFlags)
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

// RestConfig provides a mock function with given fields:
func (_m *Factory) RestConfig() (*rest.Config, error) {
	ret := _m.Called()

	var r0 *rest.Config
	if rf, ok := ret.Get(0).(func() *rest.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Config)
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

// StaticClient provides a mock function with given fields:
func (_m *Factory) StaticClient() (*kubernetes.Clientset, error) {
	ret := _m.Called()

	var r0 *kubernetes.Clientset
	if rf, ok := ret.Get(0).(func() *kubernetes.Clientset); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kubernetes.Clientset)
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