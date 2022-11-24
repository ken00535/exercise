// Code generated by MockGen. DO NOT EDIT.
// Source: ../../entity/repository_db.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	entity "shorten/internal/app/entity"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryDb is a mock of RepositoryDb interface.
type MockRepositoryDb struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryDbMockRecorder
}

// MockRepositoryDbMockRecorder is the mock recorder for MockRepositoryDb.
type MockRepositoryDbMockRecorder struct {
	mock *MockRepositoryDb
}

// NewMockRepositoryDb creates a new mock instance.
func NewMockRepositoryDb(ctrl *gomock.Controller) *MockRepositoryDb {
	mock := &MockRepositoryDb{ctrl: ctrl}
	mock.recorder = &MockRepositoryDbMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryDb) EXPECT() *MockRepositoryDbMockRecorder {
	return m.recorder
}

// GetUrl mocks base method.
func (m *MockRepositoryDb) GetUrl(ctx context.Context, shortenUrl string) (*entity.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUrl", ctx, shortenUrl)
	ret0, _ := ret[0].(*entity.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUrl indicates an expected call of GetUrl.
func (mr *MockRepositoryDbMockRecorder) GetUrl(ctx, shortenUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUrl", reflect.TypeOf((*MockRepositoryDb)(nil).GetUrl), ctx, shortenUrl)
}

// SaveShortenUrl mocks base method.
func (m *MockRepositoryDb) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveShortenUrl", ctx, url)
	ret0, _ := ret[0].(*entity.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveShortenUrl indicates an expected call of SaveShortenUrl.
func (mr *MockRepositoryDbMockRecorder) SaveShortenUrl(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveShortenUrl", reflect.TypeOf((*MockRepositoryDb)(nil).SaveShortenUrl), ctx, url)
}
