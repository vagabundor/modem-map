package queries

import (
	"errors"
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShortHandle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		handler := NewGetShortRequestHandler(mockRepo)

		mockModem := modem.ModemShort{
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
		}

		mockRepo.On("GetShort", mockModem.ID).Return(&mockModem, nil)
		mockRepo.On("RepoName", mockModem.ID.HubID).Return("TestHub", nil)

		request := GetShortRequest{ID: mockModem.ID}
		result, err := handler.Handle(request)

		expectedResult := GetShortResult{
			ID:           mockModem.ID.NetModemID,
			HubId:        mockModem.HubID,
			HubName:      "TestHub",
			ModemSn:      mockModem.ModemSn,
			NetModemName: mockModem.NetModemName,
			ActiveStatus: mockModem.ActiveStatus,
			GeoLocation: geo.ToDecimal(geo.DMS{Degrees: mockModem.Geo.LatDegrees, Minutes: mockModem.Geo.LatMinutes, Seconds: mockModem.Geo.LatSeconds, Direction: mockModem.Geo.LatSouth},
				geo.DMS{Degrees: mockModem.Geo.LongDegrees, Minutes: mockModem.Geo.LongMinutes, Seconds: mockModem.Geo.LongSeconds, Direction: mockModem.Geo.LongWest}),
		}
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting modem", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		handler := NewGetShortRequestHandler(mockRepo)

		mockID := modem.ID{
			NetModemID: 2,
			HubID:      1,
		}

		mockRepo.On("GetShort", mockID).Return((*modem.ModemShort)(nil), errors.New("error getting modem"))

		request := GetShortRequest{ID: mockID}
		_, err := handler.Handle(request)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting hub name", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		handler := NewGetShortRequestHandler(mockRepo)

		mockModem := modem.ModemShort{
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
		}

		mockRepo.On("GetShort", mockModem.ID).Return(&mockModem, nil)
		mockRepo.On("RepoName", mockModem.ID.HubID).Return("", errors.New("error getting hub name"))

		request := GetShortRequest{ID: mockModem.ID}
		_, err := handler.Handle(request)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

}
