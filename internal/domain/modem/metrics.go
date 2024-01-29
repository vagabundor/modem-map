package modem

import "modem-map/internal/pkg/geo"

// Metrics interface for modems
type Metrics interface {
	UpdateModemDetails(modem *Modem) error
	UpdateAllModems(modems []*ModemShort) error
	GetLatLong(modem *ModemShort) (geo.DD, error)
}
