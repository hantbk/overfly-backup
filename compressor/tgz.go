package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
)

// Tgz .tar.gz compressor
type Tgz struct {
}

func (ctx *Tgz) perform(model config.ModelConfig) (archivePath string, err error) {
	logger.Info("=> Compress with Tgz...")
	filePath := archiveFilePath(".tar.gz")
	_, err = helper.Exec("tar", "zcf", filePath, model.Name)
	if err == nil {
		archivePath = filePath
		return
	}
	return
}
