package mysql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"modem-map/internal/app/config"
	"modem-map/internal/domain/modem"

	"github.com/DATA-DOG/go-sqlmock"
)

// NewRepoWithDb constructor for testing purposes
func NewRepoWithDb(hubconfigs []config.HubConfig, db *gorm.DB) (Repo, error) {
	var repo Repo
	repo.names = make(map[int]string)
	if db == nil {
		return Repo{}, fmt.Errorf("failed to initialize databases")
	}
	for i, hubconfig := range hubconfigs {
		repo.dbs = append(repo.dbs, db)
		repo.names[i] = hubconfig.Hubname
	}
	return repo, nil
}

func TestRepo_GetAllShort(t *testing.T) {
	// Create the mock data
	mockModems := []*modem.ModemShort{
		{
			ID: modem.ID{
				NetModemID: 1,
				HubID:      0,
				DID:        1234567,
			},
			ModemSn:      12345,
			NetModemName: "Test Modem 1",
			ActiveStatus: 1,
			Geo: modem.Geo{
				LatDegrees:  12,
				LatMinutes:  34,
				LatSeconds:  56,
				LongDegrees: 65,
				LongMinutes: 43,
				LongSeconds: 21,
				LatSouth:    1,
				LongWest:    1,
			},
			VnoName: "Test Group 1",
		},
		{
			ID: modem.ID{
				NetModemID: 2,
				HubID:      1,
				DID:        7400,
			},
			ModemSn:      67890,
			NetModemName: "Test Modem 2",
			ActiveStatus: 1,
			Geo: modem.Geo{
				LatDegrees:  23,
				LatMinutes:  45,
				LatSeconds:  12,
				LongDegrees: 54,
				LongMinutes: 32,
				LongSeconds: 10,
				LatSouth:    0,
				LongWest:    1,
			},
			VnoName: "Test Group 2",
		},
	}

	// Create the test cases
	testCases := []struct {
		name          string
		dbConfigs     []config.HubConfig
		mockSetupFunc func(sqlmock.Sqlmock)
		expected      []*modem.ModemShort
		expectedErr   error
	}{
		{
			name: "Success",
			dbConfigs: []config.HubConfig{
				{Hubname: "test1", Dsn: ""},
				{Hubname: "test2", Dsn: ""},
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				// Add expectation for the "SELECT VERSION()" query
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))

				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, vno.Name " +
					"FROM NetModem AS nm " +
					"LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID " +
					"LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID " +
					"LEFT JOIN VNOGroupOwnedResource vnoRes ON nm.DID = vnoRes.ResourceId " +
					"LEFT JOIN VNOGroup vno ON vnoRes.GroupId = vno.ID").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus", "LatDegrees", "LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest", "Name"}).
						AddRow(mockModems[0].ID.NetModemID, mockModems[0].ID.DID, mockModems[0].ModemSn, mockModems[0].NetModemName, mockModems[0].ActiveStatus,
							mockModems[0].Geo.LatDegrees, mockModems[0].Geo.LatMinutes, mockModems[0].Geo.LatSeconds,
							mockModems[0].Geo.LongDegrees, mockModems[0].Geo.LongMinutes, mockModems[0].Geo.LongSeconds,
							mockModems[0].Geo.LatSouth, mockModems[0].Geo.LongWest, mockModems[0].VnoName))

				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, vno.Name " +
					"FROM NetModem AS nm " +
					"LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID " +
					"LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID " +
					"LEFT JOIN VNOGroupOwnedResource vnoRes ON nm.DID = vnoRes.ResourceId " +
					"LEFT JOIN VNOGroup vno ON vnoRes.GroupId = vno.ID").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus", "LatDegrees", "LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest", "Name"}).
						AddRow(mockModems[1].ID.NetModemID, mockModems[1].ID.DID, mockModems[1].ModemSn, mockModems[1].NetModemName, mockModems[1].ActiveStatus,
							mockModems[1].Geo.LatDegrees, mockModems[1].Geo.LatMinutes, mockModems[1].Geo.LatSeconds,
							mockModems[1].Geo.LongDegrees, mockModems[1].Geo.LongMinutes, mockModems[1].Geo.LongSeconds,
							mockModems[1].Geo.LatSouth, mockModems[1].Geo.LongWest, mockModems[1].VnoName))
			},
			expected:    mockModems,
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the sqlmock database connection
			sqlDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			// Apply the mock setup function
			tc.mockSetupFunc(mock)

			// Create a GORM database connection
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			require.NoError(t, err)

			// Create a new Repo using the mocked GORM database connection
			repo, err := NewRepoWithDb(tc.dbConfigs, db)
			require.NoError(t, err)

			// Call the GetAllShort method
			modems, err := repo.GetAllShort()

			// Check the results
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, modems)
			}

			// Check if there are any unfulfilled expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_Get(t *testing.T) {
	mockModem := &modem.Modem{
		ModemShort: modem.ModemShort{
			ID: modem.ID{
				NetModemID: 1,
				HubID:      0,
				DID:        1234567,
			},
			ModemSn:      12345,
			NetModemName: "Test Modem 1",
			ActiveStatus: 1,
			Geo: modem.Geo{
				LatDegrees:  12,
				LatMinutes:  34,
				LatSeconds:  56,
				LongDegrees: 78,
				LongMinutes: 90,
				LongSeconds: 12,
				LatSouth:    1,
				LongWest:    0,
			},
		},
		Model:         1,
		Buc:           "Buc",
		Lnb:           "Lnb",
		ReflectorSize: 1.2,
	}

	testCases := []struct {
		name          string
		dbConfigs     []config.HubConfig
		mockSetupFunc func(sqlmock.Sqlmock)
		input         modem.ID
		expected      *modem.Modem
		expectedErr   error
	}{
		{
			name: "Success",
			dbConfigs: []config.HubConfig{
				{Hubname: "test1", Dsn: ""},
				{Hubname: "test2", Dsn: ""},
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				// Add expectation for the "SELECT VERSION()" query
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
				// Add expectation for the query to get modem data
				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.HwType, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, r.Size, buc.ManufacturerPartNum AS Buc, " +
					"lnb.ManufacturerPartNum AS Lnb FROM NetModem AS nm").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus", "Model", "LatDegrees",
						"LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest", "ReflectorSize", "Buc", "Lnb"}).
						AddRow(mockModem.ID.NetModemID, mockModem.ID.DID, mockModem.ModemSn, mockModem.NetModemName, mockModem.ActiveStatus, mockModem.Model,
							mockModem.Geo.LatDegrees, mockModem.Geo.LatMinutes, mockModem.Geo.LatSeconds,
							mockModem.Geo.LongDegrees, mockModem.Geo.LongMinutes, mockModem.Geo.LongSeconds,
							mockModem.Geo.LatSouth, mockModem.Geo.LongWest, mockModem.ReflectorSize, mockModem.Buc, mockModem.Lnb))
			},
			input:       modem.ID{NetModemID: 1, HubID: 0, DID: 1234567},
			expected:    mockModem,
			expectedErr: nil,
		},
		{
			name: "Not found",
			dbConfigs: []config.HubConfig{
				{Hubname: "test1", Dsn: ""},
				{Hubname: "test2", Dsn: ""},
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				// Add expectation for the "SELECT VERSION()" query
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
				// Add expectation for the query to get modem data
				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.HwType, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, r.Size, " +
					"buc.ManufacturerPartNum AS Buc, lnb.ManufacturerPartNum AS Lnb FROM NetModem AS nm").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "ModemSn", "NetModemName", "ActiveStatus", "Model", "LatDegrees",
						"LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest", "ReflectorSize", "Buc", "Lnb"}))
			},
			input:       modem.ID{NetModemID: 2, HubID: 0},
			expected:    nil,
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "HubID out of range",
			dbConfigs: []config.HubConfig{
				{Hubname: "test",
					Dsn: ""},
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
			},
			input:       modem.ID{NetModemID: 1, HubID: 100}, // HubID out of range
			expected:    nil,
			expectedErr: fmt.Errorf("no repository found with id %d", 100),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			tc.mockSetupFunc(mock)

			gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
			require.NoError(t, err)

			// Create a new Repo using the mocked GORM database connection
			repo, err := NewRepoWithDb(tc.dbConfigs, gdb)
			require.NoError(t, err)

			result, err := repo.Get(tc.input)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
			if result != nil && tc.expected != nil {
				require.Equal(t, *tc.expected, *result)
			} else if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestRepo_GetShort(t *testing.T) {
	// Create the mock data
	mockModem := &modem.ModemShort{
		ID: modem.ID{
			NetModemID: 1,
			HubID:      0,
			DID:        1234567,
		},
		ModemSn:      12345,
		NetModemName: "Test Modem 1",
		ActiveStatus: 1,
		Geo: modem.Geo{
			LatDegrees:  55,
			LatMinutes:  45,
			LatSeconds:  30,
			LongDegrees: 37,
			LongMinutes: 36,
			LongSeconds: 24,
			LatSouth:    0,
			LongWest:    0,
		},
	}

	// Create the test case
	testCases := []struct {
		name          string
		dbConfig      config.HubConfig
		mockSetupFunc func(sqlmock.Sqlmock)
		inputID       modem.ID
		expected      *modem.ModemShort
		expectedErr   error
	}{
		{
			name: "Success",
			dbConfig: config.HubConfig{
				Hubname: "test",
				Dsn:     "",
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest FROM NetModem AS nm").
					WillReturnRows(sqlmock.NewRows([]string{
						"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus",
						"LatDegrees", "LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest"}).
						AddRow(
							mockModem.ID.NetModemID, mockModem.ID.DID, mockModem.ModemSn, mockModem.NetModemName, mockModem.ActiveStatus,
							mockModem.LatDegrees, mockModem.LatMinutes, mockModem.LatSeconds,
							mockModem.LongDegrees, mockModem.LongMinutes, mockModem.LongSeconds,
							mockModem.LatSouth, mockModem.LongWest,
						))
			},
			inputID:     modem.ID{NetModemID: 1, HubID: 0, DID: 1234567},
			expected:    mockModem,
			expectedErr: nil,
		},
		{
			name: "Not found",
			dbConfig: config.HubConfig{
				Hubname: "test", Dsn: "",
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				// Add expectation for the "SELECT VERSION()" query
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
				// Add expectation for the query to get modem data
				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest " +
					"FROM NetModem AS nm").
					WillReturnRows(sqlmock.NewRows([]string{
						"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus",
						"LatDegrees", "LatMinutes", "LatSeconds", "LongDegrees", "LongMinutes", "LongSeconds", "LatSouth", "LongWest"}))
			},
			inputID:     modem.ID{NetModemID: 2, HubID: 0},
			expected:    nil,
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "HubID out of range",
			dbConfig: config.HubConfig{
				Hubname: "test",
				Dsn:     "",
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
			},
			inputID:     modem.ID{NetModemID: 1, HubID: 100}, // HubID out of range
			expected:    nil,
			expectedErr: fmt.Errorf("no repository found with id %d", 100),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the sqlmock database connection
			sqlDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			// Apply the mock setup function
			tc.mockSetupFunc(mock)

			// Create a GORM database connection
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			require.NoError(t, err)

			// Create a new Repo using the mocked GORM database connection
			repo, err := NewRepoWithDb([]config.HubConfig{tc.dbConfig}, db)
			require.NoError(t, err)

			// Call the GetShort method
			modem, err := repo.GetShort(tc.inputID)

			// Check the results
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, modem)
			}

			// Check if there are any unfulfilled expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_GetAll(t *testing.T) {
	// Create the mock data
	mockModems := []*modem.Modem{
		{
			ModemShort: modem.ModemShort{
				ID: modem.ID{
					NetModemID: 1,
					HubID:      0,
					DID:        1234567,
				},
				ModemSn:      12345,
				NetModemName: "Test Modem 1",
				ActiveStatus: 1,
				Geo: modem.Geo{
					LatDegrees:  12,
					LatMinutes:  34,
					LatSeconds:  56,
					LongDegrees: 65,
					LongMinutes: 43,
					LongSeconds: 21,
					LatSouth:    1,
					LongWest:    1,
				},
			},
			Model:         1,
			Buc:           "Buc",
			Lnb:           "Lnb",
			ReflectorSize: 1.2,
		},
		{
			ModemShort: modem.ModemShort{
				ID: modem.ID{
					NetModemID: 10,
					HubID:      1,
					DID:        74000,
				},
				ModemSn:      789,
				NetModemName: "Test Modem 2",
				ActiveStatus: 1,
				Geo: modem.Geo{
					LatDegrees:  12,
					LatMinutes:  38,
					LatSeconds:  56,
					LongDegrees: 65,
					LongMinutes: 40,
					LongSeconds: 5,
					LatSouth:    1,
					LongWest:    1,
				},
			},
			Model:         3,
			Buc:           "Buc",
			Lnb:           "Lnb",
			ReflectorSize: 2.4,
		},
	}

	// Create the test cases
	testCases := []struct {
		name          string
		dbConfigs     []config.HubConfig
		mockSetupFunc func(sqlmock.Sqlmock)
		expected      []*modem.Modem
		expectedErr   error
	}{
		{
			name: "Success",
			dbConfigs: []config.HubConfig{
				{Hubname: "test1", Dsn: ""},
				{Hubname: "test2", Dsn: ""},
			},
			mockSetupFunc: func(mock sqlmock.Sqlmock) {
				// Add expectation for the "SELECT VERSION()" query
				mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))

				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.HwType, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, " +
					"r.Size, buc.ManufacturerPartNum AS Buc, lnb.ManufacturerPartNum AS Lnb " +
					"FROM NetModem AS nm " +
					"LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID " +
					"LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID " +
					"LEFT JOIN RemoteAntenna AS ra ON nm.RemoteAntennaID = ra.RemoteAntennaID " +
					"LEFT JOIN Reflector AS r ON ra.ReflectorID = r.ReflectorID " +
					"LEFT JOIN BUC AS buc ON ra.BUCID = buc.BUCID " +
					"LEFT JOIN LNB AS lnb ON ra.LNBID = lnb.LNBID").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus", "Model",
						"LatDegrees", "LatMinutes", "LatSeconds",
						"LongDegrees", "LongMinutes", "LongSeconds",
						"LatSouth", "LongWest", "ReflectorSize", "Buc", "Lnb"}).
						AddRow(mockModems[0].ID.NetModemID, mockModems[0].ID.DID, mockModems[0].ModemSn, mockModems[0].NetModemName, mockModems[0].ActiveStatus, mockModems[0].Model,
							mockModems[0].Geo.LatDegrees, mockModems[0].Geo.LatMinutes, mockModems[0].Geo.LatSeconds,
							mockModems[0].Geo.LongDegrees, mockModems[0].Geo.LongMinutes, mockModems[0].Geo.LongSeconds,
							mockModems[0].Geo.LatSouth, mockModems[0].Geo.LongWest, mockModems[0].ReflectorSize, mockModems[0].Buc, mockModems[0].Lnb))

				mock.ExpectQuery("^SELECT nm.NetModemId, nm.DID, nm.ModemSn, nm.NetModemName, nm.ActiveStatus, nm.HwType, " +
					"gl.LatDegrees, gl.LatMinutes, gl.LatSeconds, " +
					"gl.LongDegrees, gl.LongMinutes, gl.LongSeconds, " +
					"gl.LatSouth, gl.LongWest, r.Size, " +
					"buc.ManufacturerPartNum AS Buc, lnb.ManufacturerPartNum AS Lnb " +
					"FROM NetModem AS nm " +
					"LEFT JOIN Location AS loc ON nm.LocationID = loc.LocationID " +
					"LEFT JOIN GeoLocation AS gl ON loc.GeoLocationID = gl.GeoLocationID " +
					"LEFT JOIN RemoteAntenna AS ra ON nm.RemoteAntennaID = ra.RemoteAntennaID " +
					"LEFT JOIN Reflector AS r ON ra.ReflectorID = r.ReflectorID " +
					"LEFT JOIN BUC AS buc ON ra.BUCID = buc.BUCID " +
					"LEFT JOIN LNB AS lnb ON ra.LNBID = lnb.LNBID").
					WillReturnRows(sqlmock.NewRows([]string{"NetModemId", "DID", "ModemSn", "NetModemName", "ActiveStatus",
						"Model", "LatDegrees", "LatMinutes", "LatSeconds",
						"LongDegrees", "LongMinutes", "LongSeconds",
						"LatSouth", "LongWest", "ReflectorSize", "Buc", "Lnb"}).
						AddRow(mockModems[1].ID.NetModemID, mockModems[1].ID.DID, mockModems[1].ModemSn, mockModems[1].NetModemName, mockModems[1].ActiveStatus, mockModems[1].Model,
							mockModems[1].Geo.LatDegrees, mockModems[1].Geo.LatMinutes, mockModems[1].Geo.LatSeconds,
							mockModems[1].Geo.LongDegrees, mockModems[1].Geo.LongMinutes, mockModems[1].Geo.LongSeconds,
							mockModems[1].Geo.LatSouth, mockModems[1].Geo.LongWest, mockModems[1].ReflectorSize, mockModems[1].Buc, mockModems[1].Lnb))
			},
			expected:    mockModems,
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the sqlmock database connection
			sqlDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			// Apply the mock setup function
			tc.mockSetupFunc(mock)

			// Create a GORM database connection
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			require.NoError(t, err)

			// Create a new Repo using the mocked GORM database connection
			repo, err := NewRepoWithDb(tc.dbConfigs, db)
			require.NoError(t, err)

			// Call the GetAllShort method
			modems, err := repo.GetAll()

			// Check the results
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, modems)
			}

			// Check if there are any unfulfilled expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_RepoName(t *testing.T) {
	repo := Repo{
		names: map[int]string{
			0: "test1",
			1: "test2",
		},
	}

	// Create the test cases
	testCases := []struct {
		name        string
		id          int
		expected    string
		expectedErr error
	}{
		{
			name:        "Success",
			id:          0,
			expected:    "test1",
			expectedErr: nil,
		},
		{
			name:        "NotFound",
			id:          99,
			expected:    "",
			expectedErr: fmt.Errorf("no repository found with id %d", 99),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the RepoName method
			repoName, err := repo.RepoName(tc.id)

			// Check the results
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, repoName)
			}
		})
	}
}
