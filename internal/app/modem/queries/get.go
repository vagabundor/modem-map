package queries

import (
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
)

// GetResult is the return model of Modem Query Handlers
type GetResult struct {
	ID            int
	HubID         int
	HubName       string
	ModemSn       int
	NetModemName  string
	ActiveStatus  int16
	GeoLocation   geo.DD
	Model         int16
	ReflectorSize float64
	Buc           string
	Lnb           string
	Upsnr         string
	Downsnr       string
	Temperature   string
}

// GetRequest Model of the Handler
type GetRequest struct {
	ID modem.ID
}

// GetRequestHandler provides an interface to handle a GetRequest and return a GetResult
type GetRequestHandler interface {
	Handle(query GetRequest) (GetResult, error)
}

type getRequestHandler struct {
	repo modem.Repository
}

func NewGetRequestHandler(repo modem.Repository) GetRequestHandler {
	return getRequestHandler{repo: repo}
}

// Handle Handlers the GetRequest query
func (h getRequestHandler) Handle(query GetRequest) (GetResult, error) {
	m, err := h.repo.Get(query.ID)
	if err != nil {
		return GetResult{}, err
	}

	hubname, err := h.repo.RepoName(m.HubID)
	if err != nil {
		return GetResult{}, err
	}

	var result GetResult
	if (m != modem.Modem{}) {
		geo := geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
			geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})

		result = GetResult{
			ID:            m.NetModemID,
			HubID:         m.HubID,
			HubName:       hubname,
			ModemSn:       m.ModemSn,
			NetModemName:  m.NetModemName,
			ActiveStatus:  m.ActiveStatus,
			GeoLocation:   geo,
			Model:         m.Model,
			ReflectorSize: m.ReflectorSize,
			Buc:           m.Buc,
			Lnb:           m.Lnb,
		}
	}

	return result, err
}
