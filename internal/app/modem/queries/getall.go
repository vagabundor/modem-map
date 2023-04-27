package queries

import (
	"modem-map/internal/domain/modem"
	"modem-map/internal/pkg/geo"
)

// GetAllResult is the return model of Modem Query Handlers
type GetAllResult struct {
	ID            int
	HubId         int
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

// GetAllRequestHandler provides an interface to handle a GetAllRequest and return a GetAllResult
type GetAllRequestHandler interface {
	Handle() ([]GetAllResult, error)
}

type getAllRequestHandler struct {
	repo modem.Repository
}

func NewGetAllRequestHandler(repo modem.Repository) GetAllRequestHandler {
	return getAllRequestHandler{repo: repo}
}

// Handle Handlers the GetAllRequest query
func (h getAllRequestHandler) Handle() ([]GetAllResult, error) {
	ms, err := h.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []GetAllResult
	for _, m := range ms {

		hubname, err := h.repo.RepoName(m.HubID)
		if err != nil {
			return nil, err
		}

		geo := geo.ToDecimal(geo.DMS{Degrees: m.LatDegrees, Minutes: m.LatMinutes, Seconds: m.LatSeconds, Direction: m.LatSouth},
			geo.DMS{Degrees: m.LongDegrees, Minutes: m.LongMinutes, Seconds: m.LongSeconds, Direction: m.LongWest})

		result = append(result, GetAllResult{
			ID:            m.NetModemID,
			HubId:         m.HubID,
			ModemSn:       m.ModemSn,
			HubName:       hubname,
			NetModemName:  m.NetModemName,
			ActiveStatus:  m.ActiveStatus,
			GeoLocation:   geo,
			Model:         m.Model,
			ReflectorSize: m.ReflectorSize,
			Buc:           m.Buc,
			Lnb:           m.Lnb,
		})

	}
	return result, err
}
