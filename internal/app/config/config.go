package config

import (
	"github.com/gosnmp/gosnmp"
)

// DbConfig provides a struct for database config
type HubConfig struct {
	Dsn     string     `toml:"dsn"`
	Hubname string     `toml:"name"`
	Snmp    SnmpConfig `toml:"snmp"`
}

// HttpConfig provides a struct for http server config
type HttpConfig struct {
	ListenAddr   string `toml:"listen"`
	TemplatesDir string `toml:"templates"`
}

// Snmp config provides a struct for snmp config
type SnmpConfig struct {
	Address   string             `toml:"address"`
	Port      uint16             `toml:"port"`
	Community string             `toml:"community"`
	Version   gosnmp.SnmpVersion `toml:"version"`
}

// Config provides a struct for app config
type Config struct {
	Hubs []HubConfig `toml:"hub"`
	Http HttpConfig  `toml:"http"`
}

// Service read configuration for app
type Service interface {
	GetConfig() Config
}
