// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/blazee5/imageChecker/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Job is an autogenerated mock type for the Job type
type Job struct {
	mock.Mock
}

// CreateJob provides a mock function with given fields: ctx, input
func (_m *Job) CreateJob(ctx context.Context, input domain.CreateJobRequest) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for CreateJob")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CreateJobRequest) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewJob creates a new instance of Job. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJob(t interface {
	mock.TestingT
	Cleanup(func())
}) *Job {
	mock := &Job{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
