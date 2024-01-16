package modem

// Metrics interface for modems
type Metrics interface {
	UpdateModemDetails(modem *Modem) error
	UpdateAllModems(modems []*ModemShort) error
}
