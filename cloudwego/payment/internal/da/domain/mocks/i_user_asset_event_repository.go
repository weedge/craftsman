// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	station "github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

// IUserAssetEventRepository is an autogenerated mock type for the IUserAssetEventRepository type
type IUserAssetEventRepository struct {
	mock.Mock
}

// ChangeUsersAssetTx provides a mock function with given fields: ctx, event
func (_m *IUserAssetEventRepository) ChangeUsersAssetTx(ctx context.Context, event *station.BizEventAssetChange) error {
	ret := _m.Called(ctx, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *station.BizEventAssetChange) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIUserAssetEventRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserAssetEventRepository creates a new instance of IUserAssetEventRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserAssetEventRepository(t mockConstructorTestingTNewIUserAssetEventRepository) *IUserAssetEventRepository {
	mock := &IUserAssetEventRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
