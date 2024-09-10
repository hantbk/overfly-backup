package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/logger"
	"os"
	"path"
	"time"
)

type Base interface {
	perform(model config.ModelConfig) (archivePath string, err error)
}

func archiveFilePath(ext string) string {
	return path.Join(os.TempDir(), "vts-backup", time.Now().Format("2006.01.02.15.04.05")+ext)
}

func Run(model config.ModelConfig) (archivePath string, err error) {
	//logger.Info("------------- Compressor --------------")
	var ctx Base
	switch model.CompressWith.Type {
	case "tgz":
		ctx = &Tgz{}
	default:
		ctx = &Tgz{}
	}
	// set workdir
	os.Chdir(path.Join(model.DumpPath, "../"))
	archivePath, err = ctx.perform(model)
	if err != nil {
		return
	}
	logger.Info("->", archivePath)
	//logger.Info("----------- End Compressor ------------\n")
	return
}
