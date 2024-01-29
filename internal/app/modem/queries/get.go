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
	Status        string
	GeoLocation   geo.DD
	Model         int16
	ReflectorSize float64
	Buc           string
	Lnb           string
	Upsnr         string
	Downsnr       string
	Temp          string
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
	repo    modem.Repository
	metrics modem.Metrics
}

func NewGetRequestHandler(repo modem.Repository, metrics modem.Metrics) GetRequestHandler {
	return getRequestHandler{repo: repo, metrics: metrics}
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
	if m != nil {

		ismoble := m.IsMobile != 0

		var geoloc geo.DD
		if ismoble {
			geoloc, err = h.metrics.GetLatLong(&m.ModemShort)
			if err != nil {
				return result, err
			}
		} else {
			geoloc = geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
				geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})
		}

		// Get metrics data for modem
		err := h.metrics.UpdateModemDetails(m)
		if err != nil {
			return GetResult{}, err
		}

		result = GetResult{
			ID:            m.NetModemID,
			HubID:         m.HubID,
			HubName:       hubname,
			ModemSn:       m.ModemSn,
			NetModemName:  m.NetModemName,
			ActiveStatus:  m.ActiveStatus,
			Status:        m.Status,
			GeoLocation:   geoloc,
			Model:         m.Model,
			ReflectorSize: m.ReflectorSize,
			Buc:           m.Buc,
			Lnb:           m.Lnb,
			Upsnr:         m.Upsnr,
			Downsnr:       m.Downsnr,
			Temp:          m.Temp,
		}
	}

	return result, nil
}
