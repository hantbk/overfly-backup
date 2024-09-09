package main

import (
	"github.com/hantbk/vts-backup/compressor"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
	"github.com/hantbk/vts-backup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (ctx Model) perform() {
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath)
	defer ctx.cleanup()

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(ctx.Config, *archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

}

// Cleanup model temp files
func (ctx Model) cleanup() {
	logger.Info("Cleanup temp dir...")
	helper.Exec("rm", "-rf", ctx.Config.DumpPath)
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
