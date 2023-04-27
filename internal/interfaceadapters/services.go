package interfaceadapters

import (
	"modem-map/internal/app/config"
	"modem-map/internal/domain/modem"
	"modem-map/internal/interfaceadapters/config/file"
	"modem-map/internal/interfaceadapters/mysql"
)

// Services contains the exposed services of interdace adapters
type Services struct {
	ConfigService   config.Service
	ModemRepository modem.Repository
}

// NewServices Instantiates the interface adapter services
func NewServices(configPath string) (Services, error) {
	configService, err := file.NewConfigService(configPath)
	if err != nil {
		return Services{}, err
	}

	dbs, dbnames, err := mysql.InitDatabase(configService.Config.Databases)
	if err != nil {
		return Services{}, err
	}

	services := Services{
		ConfigService:   configService,
		ModemRepository: mysql.NewRepo(dbs, dbnames),
	}
	return services, nil
}
