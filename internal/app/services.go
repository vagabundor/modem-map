package app

import (
	"modem-map/internal/app/modem/queries"
	"modem-map/internal/domain/modem"
)

// Queries Contains all available query for this app
type Queries struct {
	Get           queries.GetRequestHandler
	GetShort      queries.GetShortRequestHandler
	GetAllHandler queries.GetAllRequestHandler
	GetAllShort   queries.GetAllShortRequestHandler
}

// ModemServices Contains the gruped queries of the app
type ModemServices struct {
	Queries Queries
}

// Services contains all exposed services of the application layer
type Services struct {
	ModemServices ModemServices
}

// NewServices Bootstraps Application Layer dependencies
func NewServices(modemRepo modem.Repository, modemMetrics modem.Metrics) Services {
	return Services{
		ModemServices: ModemServices{
			Queries: Queries{
				queries.NewGetRequestHandler(modemRepo, modemMetrics),
				queries.NewGetShortRequestHandler(modemRepo),
				queries.NewGetAllRequestHandler(modemRepo),
				queries.NewGetAllShortRequestHandler(modemRepo, modemMetrics),
			},
		},
	}
}
