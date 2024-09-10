package compressor

import (
	"github.com/hantbk/vts-backup/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTar_options(t *testing.T) {
	ctx := &Tar{}
	opts := ctx.options()
	if helper.IsGnuTar {
		assert.Equal(t, "--ignore-failed-read", opts[0])
	}
}
