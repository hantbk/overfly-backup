package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelsLength(t *testing.T) {
	assert.Equal(t, Exist, true)
	assert.Len(t, Models, 3)
}

func TestModel(t *testing.T) {
	model := GetModelByName("base_test")

	assert.Equal(t, model.Name, "base_test")

	// compress_with
	assert.Equal(t, model.CompressWith.Type, "tgz")
	assert.NotNil(t, model.CompressWith.Viper)

	// store_with
	assert.Equal(t, model.StoreWith.Type, "local")
	assert.Equal(t, model.StoreWith.Viper.GetString("path"), "/Users/jason/Downloads/backup1")

	// archive
	includes := model.Archive.GetStringSlice("includes")
	assert.Len(t, includes, 4)
	assert.Contains(t, includes, "/home/ubuntu/.ssh/")
	assert.Contains(t, includes, "/etc/nginx/nginx.conf")

	excludes := model.Archive.GetStringSlice("excludes")
	assert.Len(t, excludes, 2)
	assert.Contains(t, excludes, "/home/ubuntu/.ssh/known_hosts")
}
