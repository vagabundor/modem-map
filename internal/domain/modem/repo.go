package modem

// Repository Interface for modems
type Repository interface {
	Get(id ID) (*Modem, error)
	GetShort(id ID) (*ModemShort, error)
	GetAll() ([]*Modem, error)
	GetAllShort() ([]*ModemShort, error)
	RepoName(hubid int) (string, error)
}
