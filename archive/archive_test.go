package archive

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	// with nil Archive
	model := config.ModelConfig{
		Archive: nil,
	}
	err := Run(model)
	assert.NoError(t, err)
}
