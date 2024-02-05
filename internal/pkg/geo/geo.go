package geo

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
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

// StringToDecimal converts String to DD
func StringToDecimal(input string) (DD, error) {
	var result DD

	// Coordinates with N/S, E/W
	reWithDirection := regexp.MustCompile(`LAT = ([\d\.]+)(N|S) LONG = ([\d\.]+)(E|W)`)
	matchesWithDirection := reWithDirection.FindStringSubmatch(input)

	if len(matchesWithDirection) == 5 {
		return parseCoordinates(matchesWithDirection)
	}

	// Simple coordinates
	reSimple := regexp.MustCompile(`LAT-LONG: \[([\d\.]+) ([\d\.]+)\]`)
	matchesSimple := reSimple.FindStringSubmatch(input)

	if len(matchesSimple) == 3 {
		lat, err := strconv.ParseFloat(matchesSimple[1], 64)
		if err != nil {
			return result, fmt.Errorf("invalid latitude")
		}

		long, err := strconv.ParseFloat(matchesSimple[2], 64)
		if err != nil {
			return result, fmt.Errorf("invalid longitude")
		}

		result.Lat = lat
		result.Long = long
		return result, nil
	}

	return result, fmt.Errorf("invalid format")
}

func parseCoordinates(matches []string) (DD, error) {
	var result DD
	lat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return result, fmt.Errorf("invalid latitude")
	}
	if matches[2] == "S" {
		lat = -lat
	}

	long, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return result, fmt.Errorf("invalid longitude")
	}
	if matches[4] == "W" {
		long = -long
	}

	result.Lat = lat
	result.Long = long
	return result, nil
}
