package queries

import (
	"errors"
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHandle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetRequestHandler(mockRepo, mockMetrics)

		mockModem := modem.Modem{
			ModemShort: modem.ModemShort{
				ID: modem.ID{
					NetModemID: 1,
					HubID:      1,
				},
				ModemSn:      123,
				NetModemName: "test",
				ActiveStatus: 1,
				Geo: modem.Geo{
					LatDegrees:  40,
					LatMinutes:  26,
					LatSeconds:  46,
					LongDegrees: 79,
					LongMinutes: 56,
					LongSeconds: 55,
					LatSouth:    0,
					LongWest:    1,
				},
			},
			Model:         1,
			ReflectorSize: 2.5,
		}

		mockMetrics.On("UpdateModemDetails", mock.Anything).Return(nil)

		mockRepo.On("Get", mockModem.ID).Return(&mockModem, nil)
		mockRepo.On("RepoName", mockModem.ID.HubID).Return("TestHub", nil)

		request := GetRequest{ID: mockModem.ID}
		result, err := handler.Handle(request)

		expectedResult := GetResult{
			ID:           mockModem.ID.NetModemID,
			HubID:        mockModem.HubID,
			HubName:      "TestHub",
			ModemSn:      mockModem.ModemSn,
			NetModemName: mockModem.NetModemName,
			ActiveStatus: mockModem.ActiveStatus,
			GeoLocation: geo.ToDecimal(geo.DMS{Degrees: mockModem.Geo.LatDegrees, Minutes: mockModem.Geo.LatMinutes, Seconds: mockModem.Geo.LatSeconds, Direction: mockModem.Geo.LatSouth},
				geo.DMS{Degrees: mockModem.Geo.LongDegrees, Minutes: mockModem.Geo.LongMinutes, Seconds: mockModem.Geo.LongSeconds, Direction: mockModem.Geo.LongWest}),
			Model:         mockModem.Model,
			ReflectorSize: mockModem.ReflectorSize,
		}
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockRepo.AssertExpectations(t)
		mockMetrics.AssertExpectations(t)
	})

	t.Run("error getting modem", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetRequestHandler(mockRepo, mockMetrics)

		mockID := modem.ID{
			NetModemID: 2,
			HubID:      1,
		}

		mockRepo.On("Get", mockID).Return((*modem.Modem)(nil), errors.New("error getting modem"))

		request := GetRequest{ID: mockID}
		_, err := handler.Handle(request)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockMetrics.AssertExpectations(t)
	})

	t.Run("error getting hub name", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetRequestHandler(mockRepo, mockMetrics)

		mockModem := modem.Modem{
			ModemShort: modem.ModemShort{
				ID: modem.ID{
					NetModemID: 1,
					HubID:      1,
				},
				ModemSn:      123,
				NetModemName: "test",
				ActiveStatus: 1,
				Geo: modem.Geo{
					LatDegrees:  40,
					LatMinutes:  26,
					LatSeconds:  46,
					LongDegrees: 79,
					LongMinutes: 56,
					LongSeconds: 55,
					LatSouth:    0,
					LongWest:    1,
				},
			},
			Model:         1,
			ReflectorSize: 2.5,
		}

		mockRepo.On("Get", mockModem.ID).Return(&mockModem, nil)
		mockRepo.On("RepoName", mockModem.ID.HubID).Return("", errors.New("error getting hub name"))

		request := GetRequest{ID: mockModem.ID}
		_, err := handler.Handle(request)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockMetrics.AssertExpectations(t)
	})

}
