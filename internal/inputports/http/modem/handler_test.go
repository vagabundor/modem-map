package modem

import (
	"encoding/json"
	"errors"
	"fmt"
	"modem-map/internal/app"
	"modem-map/internal/app/modem/queries"
	"modem-map/internal/domain/modem"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetRequestHandler struct {
	mock.Mock
}

func (m *MockGetRequestHandler) Handle(request queries.GetRequest) (queries.GetResult, error) {
	args := m.Called(request)
	return args.Get(0).(queries.GetResult), args.Error(1)
}

type MockGetAllShortRequestHandler struct {
	mock.Mock
}

func (m *MockGetAllShortRequestHandler) Handle() ([]queries.GetAllShortResult, error) {
	args := m.Called()
	return args.Get(0).([]queries.GetAllShortResult), args.Error(1)
}

func TestGetAllShort(t *testing.T) {
	// Create a slice of modem.ModemShort
	modems := []modem.ModemShort{
		{
			ID: modem.ID{
				NetModemID: 1,
				HubID:      0,
			},
			ModemSn:      12345,
			NetModemName: "Test Modem 1",
		},
		{
			ID: modem.ID{
				NetModemID: 2,
				HubID:      1,
			},
			ModemSn:      67890,
			NetModemName: "Test Modem 2",
		},
	}

	// Convert the modems to GetAllShortResult
	mockModems := make([]queries.GetAllShortResult, len(modems))
	for i, modem := range modems {
		mockModems[i] = queries.GetAllShortResult{
			ID:           modem.ID.NetModemID,
			HubId:        modem.ID.HubID,
			ModemSn:      modem.ModemSn,
			NetModemName: modem.NetModemName,
		}
	}

	mockGet := new(MockGetRequestHandler)
	mockGetAllShort := new(MockGetAllShortRequestHandler)
	mockGetAllShort.On("Handle").Return(mockModems, nil)

	modemServices := app.ModemServices{
		Queries: app.Queries{
			Get:         mockGet,
			GetAllShort: mockGetAllShort,
		},
	}

	// Create the handler and set up the HTTP request and response recorder
	handler := NewHandler(modemServices)
	req := httptest.NewRequest("GET", "/modems", nil)
	res := httptest.NewRecorder()

	// Call the GetAllShort function
	handler.GetAllShort(res, req)

	// Check the results
	// Check the results
	assert.Equal(t, http.StatusOK, res.Code)

	var responseModems []queries.GetAllShortResult
	err := json.NewDecoder(res.Body).Decode(&responseModems)
	assert.NoError(t, err)
	assert.Equal(t, mockModems, responseModems)

}

func TestGetById(t *testing.T) {
	// Create a modem.Modem instance
	mockModem := modem.Modem{
		ModemShort: modem.ModemShort{
			ID: modem.ID{
				NetModemID: 1,
				HubID:      2,
			},
			ModemSn:      12345,
			NetModemName: "Test Modem 1",
		},
	}

	mockGetResult := queries.GetResult{
		ID:           mockModem.ID.NetModemID,
		HubID:        mockModem.ID.HubID,
		ModemSn:      mockModem.ModemSn,
		NetModemName: mockModem.NetModemName,
		// добавьте остальные поля, если они необходимы
	}

	mockGet := new(MockGetRequestHandler)
	mockGetAllShort := new(MockGetAllShortRequestHandler)

	// Настройте ожидания для mockGet
	mockGet.On("Handle", queries.GetRequest{ID: modem.ID{NetModemID: 1, HubID: 2}}).Return(mockGetResult, nil)

	modemServices := app.ModemServices{
		Queries: app.Queries{
			Get:         mockGet,
			GetAllShort: mockGetAllShort,
		},
	}

	// Create the handler and set up the HTTP request and response recorder
	handler := NewHandler(modemServices)
	req := httptest.NewRequest("GET", "/modems/2/1", nil)
	res := httptest.NewRecorder()
	// Set up the mux router to extract URL parameters
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/modems/{%s}/{%s}", GetHubIDURLParam, GetModemIDURLParam), handler.GetById)
	router.ServeHTTP(res, req)

	// Check the results
	assert.Equal(t, http.StatusOK, res.Code)

	var responseModem queries.GetResult
	err := json.NewDecoder(res.Body).Decode(&responseModem)
	assert.NoError(t, err)
	assert.Equal(t, mockGetResult, responseModem)

	// Test not found case
	mockGet.On("Handle", queries.GetRequest{ID: modem.ID{NetModemID: 5, HubID: 2}}).Return(queries.GetResult{}, nil)
	reqNotFound := httptest.NewRequest("GET", "/modems/2/5", nil)
	resNotFound := httptest.NewRecorder()

	router.ServeHTTP(resNotFound, reqNotFound)
	assert.Equal(t, http.StatusNotFound, resNotFound.Code)

	// Test error case
	mockGet.On("Handle", queries.GetRequest{ID: modem.ID{NetModemID: 3, HubID: 2}}).Return(queries.GetResult{}, errors.New("error"))
	reqError := httptest.NewRequest("GET", "/modems/2/3", nil)
	resError := httptest.NewRecorder()

	router.ServeHTTP(resError, reqError)
	assert.Equal(t, http.StatusInternalServerError, resError.Code)
}
