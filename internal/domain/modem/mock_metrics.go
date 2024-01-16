package modem

import (
	"github.com/stretchr/testify/mock"
)

type MockMetrics struct {
	mock.Mock
}

// UpdateModemDetails provides a mock function with given fields: modem
func (_m *MockMetrics) UpdateModemDetails(modem *Modem) error {
	ret := _m.Called(modem)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Modem) error); ok {
		r0 = rf(modem)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAllModems provides a mock function with given fields: modems
func (_m *MockMetrics) UpdateAllModems(modems []*ModemShort) error {
	ret := _m.Called(modems)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*ModemShort) error); ok {
		r0 = rf(modems)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
