package modem

// Geo location
type Geo struct {
	LatDegrees  uint8   `gorm:"column:LatDegrees;type:smallint(6)"`
	LatMinutes  uint8   `gorm:"column:LatMinutes;type:smallint(6)"`
	LatSeconds  float64 `gorm:"column:LatSeconds;type:double"`
	LongDegrees uint8   `gorm:"column:LongDegrees;type:smallint(6)"`
	LongMinutes uint8   `gorm:"column:LongMinutes;type:smallint(6)"`
	LongSeconds float64 `gorm:"column:LongSeconds;type:double"`
	LatSouth    uint8   `gorm:"column:LatSouth;type:tinyint(4)"`
	LongWest    uint8   `gorm:"column:LongWest;type:tinyint(4)"`
}

type ID struct {
	NetModemID int `gorm:"column:NetModemId;type:int(11)"`
	HubID      int
	DID        int `gorm:"column:DID;type:int(10)"`
}

// Short version of Modem
type ModemShort struct {
	ID
	ModemSn      int    `gorm:"column:ModemSn;type:int(11)"`
	NetModemName string `gorm:"column:NetModemName;type:varchar(64)"`
	ActiveStatus int16  `gorm:"column:ActiveStatus;type:smallint(6)"`
	Status       string
	Geo
}

// Model that represents Modem
type Modem struct {
	ModemShort
	Upsnr         string
	Downsnr       string
	Temp          string
	Model         int16   `gorm:"column:HwType;type:smallint(6)"`
	Buc           string  `gorm:"column:Buc;type:varchar(64)"`
	Lnb           string  `gorm:"column:Lnb;type:varchar(64)"`
	ReflectorSize float64 `gorm:"column:Size;type:float)"`
}

// Gorm method to change table name for Geo
func (g *Geo) TableName() string {
	return "GeoLocation"
}
