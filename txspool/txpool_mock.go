// Code generated by MockGen. DO NOT EDIT.
// Source: ./txspool/interface.go

// Package txspool is a generated GoMock package.
package txspool

import (
	event "github.com/ethereum/go-ethereum/event"
	gomock "github.com/golang/mock/gomock"
	common "github.com/palletone/go-palletone/common"
	modules "github.com/palletone/go-palletone/dag/modules"
	reflect "reflect"
)

// MockITxPool is a mock of ITxPool interface
type MockITxPool struct {
	ctrl     *gomock.Controller
	recorder *MockITxPoolMockRecorder
}

// MockITxPoolMockRecorder is the mock recorder for MockITxPool
type MockITxPoolMockRecorder struct {
	mock *MockITxPool
}

// NewMockITxPool creates a new mock instance
func NewMockITxPool(ctrl *gomock.Controller) *MockITxPool {
	mock := &MockITxPool{ctrl: ctrl}
	mock.recorder = &MockITxPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockITxPool) EXPECT() *MockITxPoolMockRecorder {
	return m.recorder
}

// Stop mocks base method
func (m *MockITxPool) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockITxPoolMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockITxPool)(nil).Stop))
}

// AddLocal mocks base method
func (m *MockITxPool) AddLocal(tx *modules.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLocal", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLocal indicates an expected call of AddLocal
func (mr *MockITxPoolMockRecorder) AddLocal(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLocal", reflect.TypeOf((*MockITxPool)(nil).AddLocal), tx)
}

// AddRemote mocks base method
func (m *MockITxPool) AddRemote(tx *modules.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemote", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRemote indicates an expected call of AddRemote
func (mr *MockITxPoolMockRecorder) AddRemote(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemote", reflect.TypeOf((*MockITxPool)(nil).AddRemote), tx)
}

// Pending mocks base method
func (m *MockITxPool) Pending() (map[common.Hash][]*TxPoolTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pending")
	ret0, _ := ret[0].(map[common.Hash][]*TxPoolTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pending indicates an expected call of Pending
func (mr *MockITxPoolMockRecorder) Pending() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pending", reflect.TypeOf((*MockITxPool)(nil).Pending))
}

// Queued mocks base method
func (m *MockITxPool) Queued() ([]*TxPoolTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Queued")
	ret0, _ := ret[0].([]*TxPoolTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Queued indicates an expected call of Queued
func (mr *MockITxPoolMockRecorder) Queued() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Queued", reflect.TypeOf((*MockITxPool)(nil).Queued))
}

// SetPendingTxs mocks base method
func (m *MockITxPool) SetPendingTxs(unit_hash common.Hash, num uint64, txs []*modules.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPendingTxs", unit_hash, num, txs)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPendingTxs indicates an expected call of SetPendingTxs
func (mr *MockITxPoolMockRecorder) SetPendingTxs(unit_hash, num, txs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPendingTxs", reflect.TypeOf((*MockITxPool)(nil).SetPendingTxs), unit_hash, num, txs)
}

// ResetPendingTxs mocks base method
func (m *MockITxPool) ResetPendingTxs(txs []*modules.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPendingTxs", txs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPendingTxs indicates an expected call of ResetPendingTxs
func (mr *MockITxPoolMockRecorder) ResetPendingTxs(txs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPendingTxs", reflect.TypeOf((*MockITxPool)(nil).ResetPendingTxs), txs)
}

// DiscardTxs mocks base method
func (m *MockITxPool) DiscardTxs(txs []*modules.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscardTxs", txs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DiscardTxs indicates an expected call of DiscardTxs
func (mr *MockITxPoolMockRecorder) DiscardTxs(txs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscardTxs", reflect.TypeOf((*MockITxPool)(nil).DiscardTxs), txs)
}

// GetUtxoFromAll mocks base method
func (m *MockITxPool) GetUtxoFromAll(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUtxoFromAll", outpoint)
	ret0, _ := ret[0].(*modules.Utxo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUtxoFromAll indicates an expected call of GetUtxoFromAll
func (mr *MockITxPoolMockRecorder) GetUtxoFromAll(outpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUtxoFromAll", reflect.TypeOf((*MockITxPool)(nil).GetUtxoFromAll), outpoint)
}

// GetUtxoFromFree mocks base method
func (m *MockITxPool) GetUtxoFromFree(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUtxoFromFree", outpoint)
	ret0, _ := ret[0].(*modules.Utxo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUtxoFromFree indicates an expected call of GetUtxoFromFree
func (mr *MockITxPoolMockRecorder) GetUtxoFromFree(outpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUtxoFromFree", reflect.TypeOf((*MockITxPool)(nil).GetUtxoFromFree), outpoint)
}

// SubscribeTxPreEvent mocks base method
func (m *MockITxPool) SubscribeTxPreEvent(arg0 chan<- modules.TxPreEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeTxPreEvent", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeTxPreEvent indicates an expected call of SubscribeTxPreEvent
func (mr *MockITxPoolMockRecorder) SubscribeTxPreEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeTxPreEvent", reflect.TypeOf((*MockITxPool)(nil).SubscribeTxPreEvent), arg0)
}

// GetSortedTxs mocks base method
func (m *MockITxPool) GetSortedTxs() ([]*TxPoolTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSortedTxs")
	ret0, _ := ret[0].([]*TxPoolTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSortedTxs indicates an expected call of GetSortedTxs
func (mr *MockITxPoolMockRecorder) GetSortedTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSortedTxs", reflect.TypeOf((*MockITxPool)(nil).GetSortedTxs))
}

// GetTx mocks base method
func (m *MockITxPool) GetTx(hash common.Hash) (*TxPoolTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", hash)
	ret0, _ := ret[0].(*TxPoolTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx
func (mr *MockITxPoolMockRecorder) GetTx(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockITxPool)(nil).GetTx), hash)
}

// GetUnpackedTxsByAddr mocks base method
func (m *MockITxPool) GetUnpackedTxsByAddr(addr common.Address) ([]*TxPoolTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnpackedTxsByAddr", addr)
	ret0, _ := ret[0].([]*TxPoolTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnpackedTxsByAddr indicates an expected call of GetUnpackedTxsByAddr
func (mr *MockITxPoolMockRecorder) GetUnpackedTxsByAddr(addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnpackedTxsByAddr", reflect.TypeOf((*MockITxPool)(nil).GetUnpackedTxsByAddr), addr)
}

// Status mocks base method
func (m *MockITxPool) Status() (int, int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status
func (mr *MockITxPoolMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockITxPool)(nil).Status))
}

// Content mocks base method
func (m *MockITxPool) Content() (map[common.Hash]*TxPoolTransaction, map[common.Hash]*TxPoolTransaction) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Content")
	ret0, _ := ret[0].(map[common.Hash]*TxPoolTransaction)
	ret1, _ := ret[1].(map[common.Hash]*TxPoolTransaction)
	return ret0, ret1
}

// Content indicates an expected call of Content
func (mr *MockITxPoolMockRecorder) Content() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Content", reflect.TypeOf((*MockITxPool)(nil).Content))
}
