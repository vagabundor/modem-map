package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileConfigService_Read(t *testing.T) {
	configService, err := NewConfigService("testdata/config.toml")
	assert.Nil(t, err)
	conf := configService.GetConfig()

	assert.Equal(t, "user:password@tcp(127.0.0.1:3306)/database1?charset=utf8mb4&parseTime=True&loc=Local", conf.Hubs[0].Dsn)
	assert.Equal(t, "storage1", conf.Hubs[0].Hubname)
	assert.Equal(t, "user:password@tcp(127.0.0.1:5432)/database2?sslmode=disable", conf.Hubs[1].Dsn)
	assert.Equal(t, "storage2", conf.Hubs[1].Hubname)
	assert.Equal(t, "./templates/", conf.Http.TemplatesDir)
	assert.Equal(t, ":8080", conf.Http.ListenAddr)

}
