package mocks

import "github.com/stretchr/testify/mock"

// Mock para TokenProvider
type MockAuthorizationProvider struct {
	mock.Mock
}

func (m *MockAuthorizationProvider) GetAuthorization() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
