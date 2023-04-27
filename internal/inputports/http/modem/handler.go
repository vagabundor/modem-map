package modem

import (
	"encoding/json"
	"fmt"
	"modem-map/internal/app"
	"modem-map/internal/app/modem/queries"
	"modem-map/internal/domain/modem"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler Modem http request handler
type Handler struct {
	modemServices app.ModemServices
}

// NewHandler Constructor
func NewHandler(app app.ModemServices) *Handler {
	return &Handler{modemServices: app}
}

// GetAllShort Returns all shorts veiews of modem
func (m Handler) GetAllShort(w http.ResponseWriter, _ *http.Request) {
	mdms, err := m.modemServices.Queries.GetAllShort.Handle()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}

	err = json.NewEncoder(w).Encode(mdms)
	if err != nil {
		return
	}
}

// GetModemIDURLParam contains the parameter identifier to be parsed by handler
const GetModemIDURLParam = "modemId"
const GetHubIDURLParam = "hubId"

// GetById Returns the modem with the provided id
func (m Handler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modemID, err := strconv.Atoi(vars[GetModemIDURLParam])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}

	hubID, err := strconv.Atoi(vars[GetHubIDURLParam])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}

	mdm, err := m.modemServices.Queries.Get.Handle(queries.GetRequest{ID: modem.ID{NetModemID: modemID, HubID: hubID}})
	if (err == nil && mdm == queries.GetResult{}) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(mdm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
}
