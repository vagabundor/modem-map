package modem

import (
	"github.com/stretchr/testify/mock"
)

type MockModemRepository struct {
	mock.Mock
}

func (m *MockModemRepository) Get(id ID) (*Modem, error) {
	args := m.Called(id)
	return args.Get(0).(*Modem), args.Error(1)
}

func (m *MockModemRepository) GetShort(id ID) (*ModemShort, error) {
	args := m.Called(id)
	return args.Get(0).(*ModemShort), args.Error(1)
}

func (m *MockModemRepository) GetAll() ([]*Modem, error) {
	args := m.Called()
	modems := args.Get(0).([]Modem)
	modemPointers := make([]*Modem, len(modems))
	for i := range modems {
		modemPointers[i] = &modems[i]
	}
	return modemPointers, args.Error(1)
}

func (m *MockModemRepository) GetAllShort() ([]*ModemShort, error) {
	args := m.Called()
	modemShorts := args.Get(0).([]ModemShort)
	modemShortPointers := make([]*ModemShort, len(modemShorts))
	for i := range modemShorts {
		modemShortPointers[i] = &modemShorts[i]
	}
	return modemShortPointers, args.Error(1)
}

func (m *MockModemRepository) RepoName(hubID int) (string, error) {
	args := m.Called(hubID)
	return args.String(0), args.Error(1)
}
