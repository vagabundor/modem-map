package app

import (
	"testing"

	"modem-map/internal/domain/modem"

	"github.com/stretchr/testify/assert"
)

func TestNewServices(t *testing.T) {
	mockModemRepo := &modem.MockModemRepository{}

	services := NewServices(mockModemRepo)

	assert.NotNil(t, services.ModemServices.Queries.Get, "Expected Get query to be initialized, got nil")
	assert.NotNil(t, services.ModemServices.Queries.GetShort, "Expected GetShort query to be initialized, got nil")
	assert.NotNil(t, services.ModemServices.Queries.GetAllHandler, "Expected GetAllHandler query to be initialized, got nil")
	assert.NotNil(t, services.ModemServices.Queries.GetAllShort, "Expected GetAllShort query to be initialized, got nil")
}
