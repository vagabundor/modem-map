package config

// DbConfig provides a struct for database config
type DbConfig struct {
	Dsn    string `toml:"dsn"`
	Dbname string `toml:"name"`
}

// HttpConfig provides a struct for http server config
type HttpConfig struct {
	ListenAddr   string `toml:"listen"`
	TemplatesDir string `toml:"templates"`
}

// Config provides a struct for app config
type Config struct {
	Databases []DbConfig `toml:"database"`
	Http      HttpConfig `toml:"http"`
}

// Service read configuration for app
type Service interface {
	GetConfig() Config
}
