package main

import (
	"github.com/hantbk/vts-backup/compressor"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
)

func main() {
	defer cleanup()
	logger.Info("WorkDir:", config.DumpPath)

	err := compressor.Run()
	if err != nil {
		logger.Error(err)
		return
	}
}

func cleanup() {
	logger.Info("Cleaning up temp dir")
	helper.Exec("rm", "-rf", config.DumpPath)
}
