package modem

import (
	"github.com/stretchr/testify/mock"
)

type MockModemRepository struct {
	mock.Mock
}

func (m *MockModemRepository) Get(id ID) (Modem, error) {
	args := m.Called(id)
	return args.Get(0).(Modem), args.Error(1)
}

func (m *MockModemRepository) GetShort(id ID) (ModemShort, error) {
	args := m.Called(id)
	return args.Get(0).(ModemShort), args.Error(1)
}

func (m *MockModemRepository) GetAll() ([]Modem, error) {
	args := m.Called()
	return args.Get(0).([]Modem), args.Error(1)
}

func (m *MockModemRepository) GetAllShort() ([]ModemShort, error) {
	args := m.Called()
	return args.Get(0).([]ModemShort), args.Error(1)
}

func (m *MockModemRepository) RepoName(hubID int) (string, error) {
	args := m.Called(hubID)
	return args.String(0), args.Error(1)
}
