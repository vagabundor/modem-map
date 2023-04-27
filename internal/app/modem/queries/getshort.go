package queries

import (
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
)

// GetShortResult is the return model of Modem Query Handlers
type GetShortResult struct {
	ID           int
	HubId        int
	HubName      string
	ModemSn      int
	NetModemName string
	ActiveStatus int16
	GeoLocation  geo.DD
}

// GetShortRequest Model of the Handler
type GetShortRequest struct {
	ID modem.ID
}

// GetShortRequestHandler provides an interface to handle a GetShortRequest and return a GetShortResult
type GetShortRequestHandler interface {
	Handle(query GetShortRequest) (GetShortResult, error)
}

type getShortRequestHandler struct {
	repo modem.Repository
}

func NewGetShortRequestHandler(repo modem.Repository) GetShortRequestHandler {
	return getShortRequestHandler{repo: repo}
}

// Handle Handlers the GetShortRequest querymodem.Modem
func (h getShortRequestHandler) Handle(query GetShortRequest) (GetShortResult, error) {
	m, err := h.repo.GetShort(query.ID)
	if err != nil {
		return GetShortResult{}, err
	}

	hubname, err := h.repo.RepoName(m.HubID)
	if err != nil {
		return GetShortResult{}, err
	}

	var result GetShortResult
	if (m != modem.ModemShort{}) {
		geo := geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
			geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})

		result = GetShortResult{
			ID:           m.NetModemID,
			HubId:        m.HubID,
			ModemSn:      m.ModemSn,
			HubName:      hubname,
			NetModemName: m.NetModemName,
			ActiveStatus: m.ActiveStatus,
			GeoLocation:  geo,
		}
	}

	return result, err
}
