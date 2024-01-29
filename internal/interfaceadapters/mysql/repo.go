package mysql

import (
	"fmt"
	"modem-map/internal/app/config"
	"modem-map/internal/domain/modem"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Repo Implements the Repository Interface to provide mysql storages
type Repo struct {
	dbs   []*gorm.DB
	names map[int]string
}

// InitDatabase opens databases
func InitDatabase(hubconfigs []config.HubConfig) ([]*gorm.DB, map[int]string, error) {
	var dbs []*gorm.DB
	var err error
	names := make(map[int]string)

	for i, hubconfig := range hubconfigs {
		newLogger := logger.Default.LogMode(logger.Silent)
		db, err := gorm.Open(mysql.Open(hubconfig.Dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			return nil, nil, err
		}
		dbs = append(dbs, db)
		names[i] = hubconfig.Hubname
	}
	return dbs, names, err
}

// NewRepo constructor
func NewRepo(dbs []*gorm.DB, names map[int]string) Repo {
	repo := Repo{dbs: dbs, names: names}
	return repo
}

// GetAllShort Returns all short versions of modems
func (sq Repo) GetAllShort() ([]*modem.ModemShort, error) {
	var allModems []*modem.ModemShort
	for i, db := range sq.dbs {
		var modems []*modem.ModemShort
		result := db.Table("NetModem AS nm").
			Select("nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.IsMobile, "+
				"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, "+
				"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, gl.LatSouth, gl.LongWest, "+
				"vno.Name").
			Joins("LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID").
			Joins("LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID").
			Joins("LEFT JOIN VNOGroupOwnedResource vnoRes ON nm.NetModemId = vnoRes.ResourceId").
			Joins("LEFT JOIN VNOGroup vno ON vnoRes.GroupId = vno.ID").
			Where("nm.NetModemTypeId = ?", 3).
			Find(&modems)

		if result.Error != nil {
			return nil, result.Error
		}
		for _, m := range modems {
			m.HubID = i
			allModems = append(allModems, m)
		}
	}
	return allModems, nil
}

// Get Returns the Modem with provided hub id and modem id
func (sq Repo) Get(id modem.ID) (*modem.Modem, error) {
	var m *modem.Modem

	// Check if id.HubID is not out of range dbs slice
	if !(id.HubID >= 0 && id.HubID < len(sq.dbs)) {
		return nil, fmt.Errorf("no repository found with id %d", id.HubID)
	}

	db := sq.dbs[id.HubID]
	result := db.Table("NetModem AS nm").
		Select("nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.IsMobile, nm.HwType, "+
			"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, "+
			"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, gl.LatSouth, gl.LongWest, "+
			"r.Size, buc.ManufacturerPartNum AS Buc, lnb.ManufacturerPartNum AS Lnb").
		Joins("LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID").
		Joins("LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID").
		Joins("LEFT JOIN RemoteAntenna AS ra ON nm.RemoteAntennaID = ra.RemoteAntennaID").
		Joins("LEFT JOIN Reflector AS r ON ra.ReflectorID = r.ReflectorID").
		Joins("LEFT JOIN BUC AS buc ON ra.BUCID = buc.BUCID").
		Joins("LEFT JOIN LNB AS lnb ON ra.LNBID = lnb.LNBID").
		Where("nm.NetModemId = ?", id.NetModemID).
		Where("nm.NetModemTypeId = ?", 3).
		First(&m)

	if result.Error != nil {
		return nil, result.Error
	}
	// If we have modem in result we change HubID
	if result.RowsAffected > 0 {
		m.HubID = id.HubID
	}

	return m, nil
}

// GetShort Returns the short version of modem with provided hub id and modem id
func (sq Repo) GetShort(id modem.ID) (*modem.ModemShort, error) {
	var m *modem.ModemShort

	//Check if id.HubID is not out of range dbs slice
	if !(id.HubID >= 0 && id.HubID < len(sq.dbs)) {
		return nil, fmt.Errorf("no repository found with id %d", id.HubID)
	}

	db := sq.dbs[id.HubID]
	result := db.Table("NetModem AS nm").
		Select("nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, "+
			"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, "+
			"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, gl.LatSouth, gl.LongWest").
		Joins("LEFT JOIN Location ON nm.LocationID = Location.LocationID").
		Joins("LEFT JOIN GeoLocation gl ON Location.GeoLocationID = gl.GeoLocationID").
		Where("nm.NetModemId = ?", id.NetModemID).
		Where("nm.NetModemTypeId = ?", 3).
		First(&m)

	if result.Error != nil {
		return nil, result.Error
	}

	//If we have modem in result we change HubID
	if result.RowsAffected > 0 {
		m.HubID = id.HubID
	}

	return m, nil
}

// GetAll Returns all stored modems
func (sq Repo) GetAll() ([]*modem.Modem, error) {
	var allModems []*modem.Modem
	for i, db := range sq.dbs {
		var modems []*modem.Modem
		result := db.Table("NetModem AS nm").
			Select(`nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.HwType,
		gl.LatDegrees, gl.LatMinutes, gl.LatSeconds,
		gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, gl.LatSouth, gl.LongWest,
		r.Size, buc.ManufacturerPartNum AS Buc, lnb.ManufacturerPartNum AS Lnb`).
			Joins("LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID").
			Joins("LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID").
			Joins("LEFT JOIN RemoteAntenna AS ra ON nm.RemoteAntennaID = ra.RemoteAntennaID").
			Joins("LEFT JOIN Reflector AS r ON ra.ReflectorID = r.ReflectorID").
			Joins("LEFT JOIN BUC AS buc ON ra.BUCID = buc.BUCID").
			Joins("LEFT JOIN LNB AS lnb ON ra.LNBID = lnb.LNBID").
			Where("nm.NetModemTypeId = ?", 3).
			Find(&modems)

		if result.Error != nil {
			return nil, result.Error
		}

		for _, m := range modems {
			m.HubID = i
			allModems = append(allModems, m)
		}
	}
	return allModems, nil
}

// RepoName returns name of repo by id
func (sq Repo) RepoName(hubid int) (string, error) {
	name, ok := sq.names[hubid]
	if !ok {
		return "", fmt.Errorf("no repository found with id %d", hubid)
	}
	return name, nil
}
