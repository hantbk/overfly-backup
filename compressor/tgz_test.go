package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTgz(t *testing.T) {
	var ctx Base
	ctx = &Tgz{}
	model := config.ModelConfig{
		Name: "test-tar",
	}
	_, err := ctx.perform(model)
	assert.Error(t, err)
}
