package compressor

import (
	"github.com/hantbk/vts-backup/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTar_options(t *testing.T) {
	tar := &Tar{}
	opts := tar.options()
	if helper.IsGnuTar {
		assert.Equal(t, opts[0], "--ignore-failed-read")
		assert.Equal(t, opts[1], "-a")
		assert.Equal(t, opts[2], "-cf")
	} else {
		assert.Equal(t, opts[0], "-a")
		assert.Equal(t, opts[1], "-cf")
	}
}
