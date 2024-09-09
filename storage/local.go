package storage

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
)

type Local struct {
}

func (ctx *Local) perform(model config.ModelConfig, archivePath string) error {
	logger.Info("=> storage | Local")
	destPath := model.StoreWith.Viper.GetString("path")
	helper.MkdirPath(destPath)
	_, err := helper.Exec("cp", archivePath, destPath)
	if err != nil {
		return err
	}
	logger.Info("Store success", destPath)
	return nil
}
