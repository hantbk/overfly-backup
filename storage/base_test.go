package storage

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase_newBase(t *testing.T) {
	model := config.ModelConfig{}
	archivePath := "/tmp/vtsbackup/test-storage/foo.zip"
	s := newBase(model, archivePath, config.SubConfig{})

	assert.Equal(t, s.archivePath, archivePath)
	assert.Equal(t, s.model, model)
	assert.Equal(t, s.viper, model.Viper)
	assert.Equal(t, s.keep, 0)
}
