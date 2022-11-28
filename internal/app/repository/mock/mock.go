// Code generated by MockGen. DO NOT EDIT.
// Source: ../entity/repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	entity "shorten/internal/app/entity"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetUrl mocks base method.
func (m *MockRepository) GetUrl(ctx context.Context, shortenUrl string) (*entity.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUrl", ctx, shortenUrl)
	ret0, _ := ret[0].(*entity.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUrl indicates an expected call of GetUrl.
func (mr *MockRepositoryMockRecorder) GetUrl(ctx, shortenUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUrl", reflect.TypeOf((*MockRepository)(nil).GetUrl), ctx, shortenUrl)
}

// SaveShortenUrl mocks base method.
func (m *MockRepository) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveShortenUrl", ctx, url)
	ret0, _ := ret[0].(*entity.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveShortenUrl indicates an expected call of SaveShortenUrl.
func (mr *MockRepositoryMockRecorder) SaveShortenUrl(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveShortenUrl", reflect.TypeOf((*MockRepository)(nil).SaveShortenUrl), ctx, url)
}