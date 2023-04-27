package inputports

import (
	"modem-map/internal/app"
	"modem-map/internal/inputports/http"
)

// Service contains the ports services
type Services struct {
	Server *http.Server
}

// NewServices instantinates the service of input ports
func NewServices(appServices app.Services, templatesDir, listenAddr string) Services {
	return Services{Server: http.NewServer(appServices, templatesDir, listenAddr)}
}
