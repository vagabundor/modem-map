package queries

import (
	"errors"
	"testing"

	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllShortRequestHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetAllShortRequestHandler(mockRepo, mockMetrics)

		mockModems := []modem.ModemShort{
			{

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
		}

		mockMetrics.On("UpdateOnlineStatus", mock.Anything).Return(nil)

		mockRepo.On("GetAllShort").Return(mockModems, nil)
		mockRepo.On("RepoName", mockModems[0].HubID).Return("hub1", nil)

		results, err := handler.Handle()

		expectedResults := []GetAllShortResult{
			{
				ID:           mockModems[0].NetModemID,
				HubId:        mockModems[0].HubID,
				ModemSn:      mockModems[0].ModemSn,
				HubName:      "hub1",
				NetModemName: mockModems[0].NetModemName,
				ActiveStatus: mockModems[0].ActiveStatus,
				GeoLocation: geo.ToDecimal(geo.DMS{Degrees: mockModems[0].LatDegrees, Minutes: mockModems[0].LatMinutes, Seconds: mockModems[0].LatSeconds, Direction: mockModems[0].LatSouth},
					geo.DMS{Degrees: mockModems[0].LongDegrees, Minutes: mockModems[0].LongMinutes, Seconds: mockModems[0].LongSeconds, Direction: mockModems[0].LongWest}),
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting modems", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetAllShortRequestHandler(mockRepo, mockMetrics)

		mockRepo.On("GetAllShort").Return([]modem.ModemShort{}, errors.New("error getting modems"))

		_, err := handler.Handle()

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting hub name", func(t *testing.T) {
		mockRepo := new(modem.MockModemRepository)
		mockMetrics := new(modem.MockMetrics)
		handler := NewGetAllShortRequestHandler(mockRepo, mockMetrics)

		mockModems := []modem.ModemShort{
			{

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
		}
		mockRepo.On("GetAllShort").Return(mockModems, nil)
		mockRepo.On("RepoName", mockModems[0].HubID).Return("", errors.New("error getting hub name"))
		mockMetrics.On("UpdateOnlineStatus", mock.Anything).Return(nil)

		_, err := handler.Handle()

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
