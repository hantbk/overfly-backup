package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
	"os"
	"path"
	"time"
)

// Tgz .tar.gz file
type Tgz struct {
}

func (ctx *Tgz) perform(model config.ModelConfig) (resultPath *string, err error) {
	logger.Info("Compressing to .tar.gz")
	archivePath := path.Join(os.TempDir(), "vts-backup", time.Now().Format(time.RFC3339)+".tar.gz")
	os.Chdir(model.DumpPath)

	_, err = helper.Exec("tar", "zcf", archivePath, "./")
	if err == nil {
		logger.Info("->", archivePath)
		resultPath = &archivePath
		return
	}
	return
}
