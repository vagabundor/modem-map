package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToDecimal(t *testing.T) {
	testCases := []struct {
		name     string
		latDMS   DMS
		longDMS  DMS
		expected DD
	}{
		{
			name: "Example 1",
			latDMS: DMS{
				Degrees:   12,
				Minutes:   34,
				Seconds:   56.789,
				Direction: 0,
			},
			longDMS: DMS{
				Degrees:   98,
				Minutes:   45,
				Seconds:   12.345,
				Direction: 0,
			},
			expected: DD{
				Lat:  12.582441,
				Long: 98.753429,
			},
		},
		{
			name: "Example 2 - South Latitude and West Longitude",
			latDMS: DMS{
				Degrees:   51,
				Minutes:   30,
				Seconds:   0,
				Direction: 1,
			},
			longDMS: DMS{
				Degrees:   0,
				Minutes:   8,
				Seconds:   0,
				Direction: 1,
			},
			expected: DD{
				Lat:  -51.500000,
				Long: -0.133333,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToDecimal(tc.latDMS, tc.longDMS)
			assert.Equal(t, tc.expected, result)
		})
	}
}
