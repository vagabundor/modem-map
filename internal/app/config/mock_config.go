package config

import "github.com/stretchr/testify/mock"

// MockConfigService reads mock config
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) Read() (Config, error) {
	args := m.Called()
	return Config{}, args.Error(0)
}
