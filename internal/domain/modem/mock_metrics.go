package modem

import (
	"modem-map/internal/pkg/geo"

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

// GetLatLong provides a mock function
func (_m *MockMetrics) GetLatLong(modem *ModemShort) (geo.DD, error) {
	ret := _m.Called(modem)

	var r0 geo.DD
	if rf, ok := ret.Get(0).(func(*ModemShort) geo.DD); ok {
		r0 = rf(modem)
	} else {
		r0 = ret.Get(0).(geo.DD)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ModemShort) error); ok {
		r1 = rf(modem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
