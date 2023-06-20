package main

import (
	"flag"
	"log"
	"modem-map/internal/app"
	"modem-map/internal/inputports"
	"modem-map/internal/interfaceadapters"
)

func main() {
	var configPath string

	default_config := "config.toml"
	flag.StringVar(&configPath, "config", default_config, "Path to the configuration file")
	flag.StringVar(&configPath, "c", default_config, "Path to the configuration file (short version)")

	// Parse the command-line flags
	flag.Parse()

	// Load the configuration from the file specified by the command-line flag
	interfaceAdaptersServices, err := interfaceadapters.NewServices(configPath)
	if err != nil {
		log.Fatalf("Error initializing interface adapters services: %v", err)
	}

	cfg := interfaceAdaptersServices.ConfigService.GetConfig()

	appServices := app.NewServices(interfaceAdaptersServices.ModemRepository, interfaceAdaptersServices.ModemMetrics)
	inputPortServices := inputports.NewServices(appServices, cfg.Http.TemplatesDir, cfg.Http.ListenAddr)
	inputPortServices.Server.ListenAndServe()
}
