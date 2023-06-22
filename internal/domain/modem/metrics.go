package modem

// Metrics interface for modems
type Metrics interface {
	UpdateMetrics(modem *Modem) error
	UpdateStatus(modems []*ModemShort) error
}
