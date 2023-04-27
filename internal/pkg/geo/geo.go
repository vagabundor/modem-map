package geo

import (
	"math"
)

// DMS is Degrees Minutes Seconds coordinates
type DMS struct {
	Degrees uint8
	Minutes uint8
	Seconds float64
	//Direction = 1 if South on the Latitude coordinate or West on the Longitude coordinate
	Direction uint8
}

// DD is decimal degrees coordianates
type DD struct {
	Lat  float64
	Long float64
}

// ToDecimal converts Degree, Minute, Seconds coordinates to Decimal Degreees, we round float result to 6 decimal places, it ensures an accuracy of about 10cm.
func ToDecimal(latdms DMS, longdms DMS) DD {
	lat := math.Round((float64(latdms.Degrees)+float64(latdms.Minutes)/60+latdms.Seconds/3600)*1000000) / 1000000
	if latdms.Direction == uint8(1) {
		lat = -lat
	}
	long := math.Round((float64(longdms.Degrees)+float64(longdms.Minutes)/60+longdms.Seconds/3600)*1000000) / 1000000
	if longdms.Direction == uint8(1) {
		long = -long
	}
	return DD{Lat: lat, Long: long}
}
