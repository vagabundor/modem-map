package file

import (
	"testing"

	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestFileConfigService_Read(t *testing.T) {
	configService, err := NewConfigService("testdata/config.toml")
	assert.Nil(t, err)
	conf := configService.GetConfig()

	assert.Equal(t, "user:password@tcp(127.0.0.1:3306)/database1?charset=utf8mb4&parseTime=True&loc=Local", conf.Hubs[0].Dsn)
	assert.Equal(t, "hub1", conf.Hubs[0].Hubname)
	assert.Equal(t, "127.0.0.1", conf.Hubs[0].Snmp.Address)
	assert.Equal(t, uint16(161), conf.Hubs[0].Snmp.Port)
	assert.Equal(t, "public", conf.Hubs[0].Snmp.Community)
	assert.Equal(t, gosnmp.SnmpVersion(0x1), conf.Hubs[0].Snmp.Version)

	assert.Equal(t, "user:password@tcp(10.0.0.1:5432)/database2?sslmode=disable", conf.Hubs[1].Dsn)
	assert.Equal(t, "hub2", conf.Hubs[1].Hubname)
	assert.Equal(t, "10.0.0.1", conf.Hubs[1].Snmp.Address)
	assert.Equal(t, uint16(161), conf.Hubs[1].Snmp.Port)
	assert.Equal(t, "public", conf.Hubs[1].Snmp.Community)
	assert.Equal(t, gosnmp.SnmpVersion(0x1), conf.Hubs[1].Snmp.Version)

	assert.Equal(t, "./templates/", conf.Http.TemplatesDir)
	assert.Equal(t, ":8080", conf.Http.ListenAddr)
}
