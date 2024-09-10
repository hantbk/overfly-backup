package model

import (
	"github.com/hantbk/vts-backup/archive"
	"github.com/hantbk/vts-backup/compressor"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/encryptor"
	"github.com/hantbk/vts-backup/logger"
	"github.com/hantbk/vts-backup/storage"
	"os"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (ctx Model) Perform() {
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")

	defer func() {
		if r := recover(); r != nil {
			ctx.cleanup()
		}

		ctx.cleanup()
	}()

	if ctx.Config.Archive != nil {
		logger.Info("------------- Archives -------------")
		err := archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info("------------- Archives -------------\n")
	}

	//logger.Info("------------ Compressor -------------")
	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}
	//logger.Info("------------ Compressor -------------\n")

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("------------- Storage --------------")
	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("------------- Storage --------------\n")
}

// Cleanup model temp files
func (ctx Model) cleanup() {
	logger.Info("Cleanup temp: " + ctx.Config.TempPath + "/\n")
	err := os.RemoveAll(ctx.Config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+ctx.Config.TempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
