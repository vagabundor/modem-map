package interfaceadapters

import (
	"modem-map/internal/app/config"
	"modem-map/internal/domain/modem"
	"modem-map/internal/interfaceadapters/config/file"
	"modem-map/internal/interfaceadapters/mysql"
	"modem-map/internal/interfaceadapters/snmp"
)

// Services contains the exposed services of interdace adapters
type Services struct {
	ConfigService   config.Service
	ModemRepository modem.Repository
	ModemMetrics    modem.Metrics
}

// NewServices Instantiates the interface adapter services
func NewServices(configPath string) (Services, error) {
	configService, err := file.NewConfigService(configPath)
	if err != nil {
		return Services{}, err
	}

	dbs, dbnames, err := mysql.InitDatabase(configService.Config.Hubs)
	if err != nil {
		return Services{}, err
	}

	snmps, err := snmp.InitSNMP(configService.Config.Hubs)
	if err != nil {
		return Services{}, err
	}

	services := Services{
		ConfigService:   configService,
		ModemRepository: mysql.NewRepo(dbs, dbnames),
		ModemMetrics:    snmp.NewMetrics(snmps),
	}
	return services, nil
}
