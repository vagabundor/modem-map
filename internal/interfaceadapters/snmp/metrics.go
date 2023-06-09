package snmp

import (
	"fmt"
	"modem-map/internal/app/config"
	"modem-map/internal/domain/modem"
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Metrics Implements the Metrics Interface to provide live data of modems
type Metrics struct {
	snmps []*gosnmp.GoSNMP
	oids  map[string]string
}

func InitSNMP(hubconfigs []config.HubConfig) ([]*gosnmp.GoSNMP, error) {

	var snmps []*gosnmp.GoSNMP

	for _, hubconfig := range hubconfigs {

		snmpClient := &gosnmp.GoSNMP{
			Target:    hubconfig.Snmp.Address,
			Port:      hubconfig.Snmp.Port,
			Community: hubconfig.Snmp.Community,
			Version:   hubconfig.Snmp.Version,
			Timeout:   time.Duration(30) * time.Second,
		}

		err := snmpClient.Connect()
		if err != nil {
			err = fmt.Errorf("Connection: %w", err)
			return nil, err
		}

		snmps = append(snmps, snmpClient)
	}

	return snmps, nil
}

func (m *Metrics) Close() error {
	for _, snmpClient := range m.snmps {
		err := snmpClient.Conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMetrics(snmps []*gosnmp.GoSNMP) Metrics {
	metrics := Metrics{
		snmps: snmps,
		oids: map[string]string{
			"Upsnr":   ".1.3.6.1.4.1.13732.1.4.4.1.2.",
			"Downsnr": ".1.3.6.1.4.1.13732.1.4.3.1.2.",
			"Temp":    ".1.3.6.1.4.1.13732.1.4.3.1.8.",
			"Status":  ".1.3.6.1.4.1.13732.1.1.1.1.15.",
		},
	}
	return metrics
}

func (m Metrics) UpdateMetrics(modem *modem.Modem) error {

	oidslist := make([]string, 0, len(m.oids))

	for _, oid := range m.oids {
		oidslist = append(oidslist, oid+strconv.Itoa(modem.DID))
	}

	i := modem.HubID
	resp, err := m.snmps[i].Get(oidslist)
	if err != nil {
		err = fmt.Errorf("Requesting snmp: %w", err)
		return err
	}

	strDID := strconv.Itoa(modem.DID)
	for _, variable := range resp.Variables {
		switch variable.Name {
		case m.oids["Upsnr"] + strDID:
			if val, ok := variable.Value.([]uint8); ok {
				modem.Upsnr = string(val)
			}
		case m.oids["Downsnr"] + strDID:
			if val, ok := variable.Value.([]uint8); ok {
				modem.Downsnr = string(val)
			}
		case m.oids["Temp"] + strDID:
			if val, ok := variable.Value.([]uint8); ok {
				modem.Temp = string(val)
			}
		case m.oids["Status"] + strDID:
			if variable.Value != nil {
				if val, ok := variable.Value.([]uint8); ok {
					modem.Status = string(val)
				}
			}
		}
	}
	return nil
}

func (m Metrics) UpdateStatus(modems []*modem.ModemShort) error {
	for _, modem := range modems {
		i := modem.HubID
		resp, err := m.snmps[i].Get([]string{m.oids["Status"] + strconv.Itoa(modem.DID)})
		if err != nil {
			err = fmt.Errorf("Requesting snmp: %w", err)
			return err
		}

		strDID := strconv.Itoa(modem.DID)
		for _, variable := range resp.Variables {
			switch variable.Name {
			case m.oids["Status"] + strDID:
				if variable.Value != nil {
					if val, ok := variable.Value.([]uint8); ok {
						modem.Status = string(val)
					}
				}
			}
		}
	}
	return nil
}
