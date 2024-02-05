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

func TestStringToDecimal(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  DD
		expectErr bool
	}{
		{
			name:  "Valid Coordinates with Direction",
			input: "LAT-LONG : [LAT = 68.96984N LONG = 33.05111E]",
			expected: DD{
				Lat:  68.96984,
				Long: 33.05111,
			},
			expectErr: false,
		},
		{
			name:      "Invalid Format",
			input:     "LATITUDE = 68.96968N LONGITUDE = 33.05082E",
			expected:  DD{},
			expectErr: true,
		},
		{
			name:  "Negative Coordinates - South and West",
			input: "LAT = 10.12345S LONG = 20.67890W",
			expected: DD{
				Lat:  -10.12345,
				Long: -20.67890,
			},
			expectErr: false,
		},
		{
			name:  "Valid Simple Coordinates",
			input: "LAT-LONG: [64.522 40.526]",
			expected: DD{
				Lat:  64.522,
				Long: 40.526,
			},
			expectErr: false,
		},
		{
			name:      "Coordinates Not Found",
			input:     "LAT-LONG:[ Not Found ]",
			expected:  DD{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToDecimal(tc.input)

			if tc.expectErr {
				assert.Error(t, err, "Expected an error for case: %v", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for case: %v", tc.name)
				assert.Equal(t, tc.expected, result, "Expected and actual values differ for case: %v", tc.name)
			}
		})
	}
}
