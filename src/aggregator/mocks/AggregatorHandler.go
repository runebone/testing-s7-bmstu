// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// AggregatorHandler is an autogenerated mock type for the AggregatorHandler type
type AggregatorHandler struct {
	mock.Mock
}

// CreateBoard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// CreateCard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// CreateColumn provides a mock function with given fields: w, r
func (_m *AggregatorHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DeleteBoard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DeleteCard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DeleteColumn provides a mock function with given fields: w, r
func (_m *AggregatorHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetBoard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetBoards provides a mock function with given fields: w, r
func (_m *AggregatorHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetCard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetColumn provides a mock function with given fields: w, r
func (_m *AggregatorHandler) GetColumn(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetStats provides a mock function with given fields: w, r
func (_m *AggregatorHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// Login provides a mock function with given fields: w, r
func (_m *AggregatorHandler) Login(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// Logout provides a mock function with given fields: w, r
func (_m *AggregatorHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// Refresh provides a mock function with given fields: w, r
func (_m *AggregatorHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// Register provides a mock function with given fields: w, r
func (_m *AggregatorHandler) Register(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// UpdateBoard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// UpdateCard provides a mock function with given fields: w, r
func (_m *AggregatorHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// UpdateColumn provides a mock function with given fields: w, r
func (_m *AggregatorHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// Validate provides a mock function with given fields: w, r
func (_m *AggregatorHandler) Validate(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewAggregatorHandler creates a new instance of AggregatorHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAggregatorHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *AggregatorHandler {
	mock := &AggregatorHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}