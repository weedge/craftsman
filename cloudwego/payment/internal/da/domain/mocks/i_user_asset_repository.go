// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

// IUserAssetRepository is an autogenerated mock type for the IUserAssetRepository type
type IUserAssetRepository struct {
	mock.Mock
}

// GetUserAssets provides a mock function with given fields: ctx, userIds
func (_m *IUserAssetRepository) GetUserAssets(ctx context.Context, userIds []int64) ([]model.UserAsset, error) {
	ret := _m.Called(ctx, userIds)

	var r0 []model.UserAsset
	if rf, ok := ret.Get(0).(func(context.Context, []int64) []model.UserAsset); ok {
		r0 = rf(ctx, userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserAsset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int64) error); ok {
		r1 = rf(ctx, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIUserAssetRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserAssetRepository creates a new instance of IUserAssetRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserAssetRepository(t mockConstructorTestingTNewIUserAssetRepository) *IUserAssetRepository {
	mock := &IUserAssetRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
