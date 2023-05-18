// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	app "github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

// Mockstorager is a mock of storager interface.
type Mockstorager struct {
	ctrl     *gomock.Controller
	recorder *MockstoragerMockRecorder
}

// MockstoragerMockRecorder is the mock recorder for Mockstorager.
type MockstoragerMockRecorder struct {
	mock *Mockstorager
}

// NewMockstorager creates a new mock instance.
func NewMockstorager(ctrl *gomock.Controller) *Mockstorager {
	mock := &Mockstorager{ctrl: ctrl}
	mock.recorder = &MockstoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorager) EXPECT() *MockstoragerMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *Mockstorager) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockstoragerMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*Mockstorager)(nil).Ping), ctx)
}

// MockorderRepository is a mock of orderRepository interface.
type MockorderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockorderRepositoryMockRecorder
}

// MockorderRepositoryMockRecorder is the mock recorder for MockorderRepository.
type MockorderRepositoryMockRecorder struct {
	mock *MockorderRepository
}

// NewMockorderRepository creates a new mock instance.
func NewMockorderRepository(ctrl *gomock.Controller) *MockorderRepository {
	mock := &MockorderRepository{ctrl: ctrl}
	mock.recorder = &MockorderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockorderRepository) EXPECT() *MockorderRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockorderRepository) Add(ctx context.Context, orderID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockorderRepositoryMockRecorder) Add(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockorderRepository)(nil).Add), ctx, orderID)
}

// Delete mocks base method.
func (m *MockorderRepository) Delete(ctx context.Context, orderID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockorderRepositoryMockRecorder) Delete(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockorderRepository)(nil).Delete), ctx, orderID)
}

// FindByID mocks base method.
func (m *MockorderRepository) FindByID(ctx context.Context, orderID int64) (app.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, orderID)
	ret0, _ := ret[0].(app.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockorderRepositoryMockRecorder) FindByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockorderRepository)(nil).FindByID), ctx, orderID)
}

// FindRegistered mocks base method.
func (m *MockorderRepository) FindRegistered(ctx context.Context) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRegistered", ctx)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRegistered indicates an expected call of FindRegistered.
func (mr *MockorderRepositoryMockRecorder) FindRegistered(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRegistered", reflect.TypeOf((*MockorderRepository)(nil).FindRegistered), ctx)
}

// Update mocks base method.
func (m *MockorderRepository) Update(ctx context.Context, order app.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockorderRepositoryMockRecorder) Update(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockorderRepository)(nil).Update), ctx, order)
}

// UpdateStatus mocks base method.
func (m *MockorderRepository) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, orderID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockorderRepositoryMockRecorder) UpdateStatus(ctx, orderID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockorderRepository)(nil).UpdateStatus), ctx, orderID, status)
}

// MockproductRepository is a mock of productRepository interface.
type MockproductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockproductRepositoryMockRecorder
}

// MockproductRepositoryMockRecorder is the mock recorder for MockproductRepository.
type MockproductRepositoryMockRecorder struct {
	mock *MockproductRepository
}

// NewMockproductRepository creates a new mock instance.
func NewMockproductRepository(ctrl *gomock.Controller) *MockproductRepository {
	mock := &MockproductRepository{ctrl: ctrl}
	mock.recorder = &MockproductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproductRepository) EXPECT() *MockproductRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockproductRepository) Add(ctx context.Context, orderID int64, product []app.OrderProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, orderID, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockproductRepositoryMockRecorder) Add(ctx, orderID, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockproductRepository)(nil).Add), ctx, orderID, product)
}

// FindByOrderID mocks base method.
func (m *MockproductRepository) FindByOrderID(ctx context.Context, orderID int64) (app.OrderGoods, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByOrderID", ctx, orderID)
	ret0, _ := ret[0].(app.OrderGoods)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByOrderID indicates an expected call of FindByOrderID.
func (mr *MockproductRepositoryMockRecorder) FindByOrderID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByOrderID", reflect.TypeOf((*MockproductRepository)(nil).FindByOrderID), ctx, orderID)
}

// MockrewardRepository is a mock of rewardRepository interface.
type MockrewardRepository struct {
	ctrl     *gomock.Controller
	recorder *MockrewardRepositoryMockRecorder
}

// MockrewardRepositoryMockRecorder is the mock recorder for MockrewardRepository.
type MockrewardRepositoryMockRecorder struct {
	mock *MockrewardRepository
}

// NewMockrewardRepository creates a new mock instance.
func NewMockrewardRepository(ctrl *gomock.Controller) *MockrewardRepository {
	mock := &MockrewardRepository{ctrl: ctrl}
	mock.recorder = &MockrewardRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrewardRepository) EXPECT() *MockrewardRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockrewardRepository) Add(ctx context.Context, reward app.Reward) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, reward)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockrewardRepositoryMockRecorder) Add(ctx, reward interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockrewardRepository)(nil).Add), ctx, reward)
}

// Find mocks base method.
func (m *MockrewardRepository) Find(ctx context.Context, description string) (app.Reward, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, description)
	ret0, _ := ret[0].(app.Reward)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockrewardRepositoryMockRecorder) Find(ctx, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockrewardRepository)(nil).Find), ctx, description)
}

// FindByID mocks base method.
func (m *MockrewardRepository) FindByID(ctx context.Context, ID string) (app.Reward, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, ID)
	ret0, _ := ret[0].(app.Reward)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockrewardRepositoryMockRecorder) FindByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockrewardRepository)(nil).FindByID), ctx, ID)
}