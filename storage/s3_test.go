package storage

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type serviceInfo struct {
	name, endpoint, region string
}

func Test_providerName(t *testing.T) {
	var cases = map[string]serviceInfo{
		"s3": {"AWS S3", "", "us-east-1"},
	}

	base := newBase(config.ModelConfig{}, "test")
	base.viper = viper.New()

	for service, info := range cases {
		s := &S3{Base: base, Service: service}
		s.init()

		assert.Equal(t, info.name, s.providerName(), "providerName for "+service)
		assert.Equal(t, info.endpoint, *s.defaultEndpoint(), "defaultEndpoint for "+service)
		assert.Equal(t, info.region, s.defaultRegion(), "defaultRegion for "+service)

		assert.Equal(t, info.region, s.viper.GetString("region"))
		assert.Equal(t, info.endpoint, s.viper.GetString("endpoint"))
		assert.Equal(t, "3", s.viper.GetString("max_retries"))
		assert.Equal(t, "300", s.viper.GetString("timeout"))
	}

}
