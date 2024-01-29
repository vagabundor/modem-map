package queries

import (
	"errors"
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
)

// GetAllShortResult is the return model of Modem Query Handlers
type GetAllShortResult struct {
	ID           int
	HubId        int
	HubName      string
	ModemSn      int
	NetModemName string
	ActiveStatus int16
	Status       string
	VnoName      string
	IsMobile     bool
	GeoLocation  geo.DD
}

// GetAllShortRequestHandler provides an interface to handle a GetAllShortRequest and return a GetAllShortResult
type GetAllShortRequestHandler interface {
	Handle() ([]GetAllShortResult, error)
}

type getAllShortRequestHandler struct {
	repo    modem.Repository
	metrics modem.Metrics
}

func NewGetAllShortRequestHandler(repo modem.Repository, metrics modem.Metrics) GetAllShortRequestHandler {
	return getAllShortRequestHandler{repo: repo, metrics: metrics}
}

// Handle Handlers the GetAllShortRequest query
func (h getAllShortRequestHandler) Handle() ([]GetAllShortResult, error) {
	ms, err := h.repo.GetAllShort()
	if err != nil {
		return nil, err
	}

	if ms == nil {
		return nil, errors.New("modem query: GetAllShortRequest result is nil")
	}

	// Get metrics data for modem
	err = h.metrics.UpdateAllModems(ms)
	if err != nil {
		return nil, err
	}

	var result []GetAllShortResult
	for _, m := range ms {

		hubname, err := h.repo.RepoName(m.HubID)
		if err != nil {
			return nil, err
		}

		ismoble := m.IsMobile != 0

		var geoloc geo.DD
		if ismoble {
			geoloc, err = h.metrics.GetLatLong(m)
			if err != nil {
				return nil, err
			}
		} else {
			geoloc = geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
				geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})
		}

		result = append(result, GetAllShortResult{
			ID:           m.NetModemID,
			HubId:        m.HubID,
			ModemSn:      m.ModemSn,
			HubName:      hubname,
			NetModemName: m.NetModemName,
			ActiveStatus: m.ActiveStatus,
			Status:       m.Status,
			VnoName:      m.VnoName,
			IsMobile:     ismoble,
			GeoLocation:  geoloc,
		})

	}
	return result, nil
}
