package compressor

import (
	"github.com/hantbk/vts-backup/helper"
	"os/exec"
)

// Tgz .tar.gz compressor
type Tgz struct {
	Base
}

func (ctx *Tgz) perform() (archivePath string, err error) {
	filePath := ctx.archiveFilePath(".tar.gz")

	opts := ctx.options()
	opts = append(opts, filePath)
	opts = append(opts, ctx.name)

	_, err = helper.Exec("tar", opts...)
	if err == nil {
		archivePath = filePath
		return
	}
	return
}

func (ctx *Tgz) options() (opts []string) {
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}

	path, err := exec.LookPath("pigz")
	if err == nil {
		opts = append(opts, "--use-compress-program", path, "-cf")
	} else {
		opts = append(opts, "-zcf")
	}

	return
}
