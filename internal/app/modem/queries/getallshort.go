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
	GeoLocation  geo.DD
}

// GetAllShortRequestHandler provides an interface to handle a GetAllShortRequest and return a GetAllShortResult
type GetAllShortRequestHandler interface {
	Handle() ([]GetAllShortResult, error)
}

type getAllShortRequestHandler struct {
	repo modem.Repository
}

func NewGetAllShortRequestHandler(repo modem.Repository) GetAllShortRequestHandler {
	return getAllShortRequestHandler{repo: repo}
}

// Handle Handlers the GetAllShortRequest query
func (h getAllShortRequestHandler) Handle() ([]GetAllShortResult, error) {
	ms, err := h.repo.GetAllShort()
	if err != nil {
		return nil, err
	}

	var result []GetAllShortResult
	if ms == nil {
		return nil, errors.New("modem query: GetAllShortRequest result is nil")
	}
	for _, m := range ms {

		hubname, err := h.repo.RepoName(m.HubID)
		if err != nil {
			return nil, err
		}

		geo := geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
			geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})

		result = append(result, GetAllShortResult{
			ID:           m.NetModemID,
			HubId:        m.HubID,
			ModemSn:      m.ModemSn,
			HubName:      hubname,
			NetModemName: m.NetModemName,
			ActiveStatus: m.ActiveStatus,
			GeoLocation:  geo,
		})

	}
	return result, err
}
