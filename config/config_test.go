package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	Init("../vtsbackup_test.yml")
}

func TestModelsLength(t *testing.T) {
	assert.Equal(t, Exist, true)
	assert.Len(t, Models, 5)
}

func TestModel(t *testing.T) {
	model := GetModelByName("base_test")

	assert.Equal(t, model.Name, "base_test")

	// compress_with
	assert.Equal(t, model.CompressWith.Type, "tgz")
	assert.NotNil(t, model.CompressWith.Viper)

	// encrypt_with
	assert.Equal(t, model.EncryptWith.Type, "openssl")
	assert.NotNil(t, model.EncryptWith.Viper)

	assert.Equal(t, model.Storages["local"].Type, "local")
	assert.Equal(t, model.Storages["local"].Viper.GetString("path"), "/Users/hant/Downloads/backup1")

	assert.Equal(t, model.Storages["scp"].Type, "scp")
	assert.Equal(t, model.Storages["scp"].Viper.GetString("host"), "your-host.com")

	// archive
	includes := model.Archive.GetStringSlice("includes")
	assert.Len(t, includes, 4)
	assert.Contains(t, includes, "/home/ubuntu/.ssh/")
	assert.Contains(t, includes, "/etc/nginx/nginx.conf")

	excludes := model.Archive.GetStringSlice("excludes")
	assert.Len(t, excludes, 2)
	assert.Contains(t, excludes, "/home/ubuntu/.ssh/known_hosts")
}
