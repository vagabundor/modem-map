package http

import (
	"log"
	"modem-map/internal/app"
	"modem-map/internal/inputports/http/geomap"
	"modem-map/internal/inputports/http/modem"
	"net/http"

	"github.com/gorilla/mux"
)

// Server Represents the http server running for this service
type Server struct {
	appServices  app.Services
	router       *mux.Router
	templatesDir string
	listenAddr   string
}

// NewServer HTTP Server constructor
func NewServer(appServices app.Services, templatesDir, listenAddr string) *Server {
	httpServer := &Server{appServices: appServices, templatesDir: templatesDir, listenAddr: listenAddr}
	httpServer.router = mux.NewRouter()
	httpServer.AddModemHTTPRoutes()
	httpServer.AddGeomapHttpRoutes()
	httpServer.AddStaticRoutes()
	http.Handle("/", httpServer.router)

	return httpServer
}

// AddCragHTTPRoutes registers crag route handlers
func (httpServer *Server) AddModemHTTPRoutes() {
	const HTTPRoutePath = "/modems"

	//Queries
	httpServer.router.HandleFunc(HTTPRoutePath+"/{"+modem.GetHubIDURLParam+"}"+"/{"+modem.GetModemIDURLParam+"}", modem.NewHandler(httpServer.appServices.ModemServices).GetById).Methods("GET")
	httpServer.router.HandleFunc(HTTPRoutePath, modem.NewHandler(httpServer.appServices.ModemServices).GetAllShort).Methods("GET")

}

func (httpServer *Server) AddGeomapHttpRoutes() {
	const HTTPRoutePath = "/map"

	//Pages
	httpServer.router.HandleFunc(HTTPRoutePath, geomap.NewHandler(httpServer.templatesDir).Handle).Methods("GET")
}

func (httpServer *Server) AddStaticRoutes() {
	// Serve static files from the ./static directory
	httpServer.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// ListenAndServe Starts listening for requests
func (httpServer *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(httpServer.listenAddr, nil))
}
