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

func (ctx *Tgz) perform() error {
	logger.Info("Compressing to .tar.gz")
	archivePath := path.Join(os.TempDir(), "vts-backup", time.Now().Format(time.RFC3339)+".tar.gz")
	err := os.Chdir(config.DumpPath)
	if err != nil {
		return err
	}
	_, err = helper.Exec("tar", "zcf", archivePath, "./")
	if err == nil {
		logger.Info("->", archivePath)
	}
	return err
}
