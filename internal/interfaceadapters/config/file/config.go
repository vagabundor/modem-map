package file

import (
	"io/ioutil"

	"modem-map/internal/app/config"

	"github.com/BurntSushi/toml"
)

// ConfigService service for reading config from File
type ConfigService struct {
	FilePath string
	Config   config.Config
}

func NewConfigService(filePath string) (*ConfigService, error) {
	configService := &ConfigService{FilePath: filePath}

	// Read file to memory
	data, err := ioutil.ReadFile(configService.FilePath)
	if err != nil {
		return nil, err
	}

	// TOML parsing
	if _, err := toml.Decode(string(data), &configService.Config); err != nil {
		return nil, err
	}

	return configService, nil
}

// GetConfig returns the stored configuration
func (s ConfigService) GetConfig() config.Config {
	return s.Config
}
