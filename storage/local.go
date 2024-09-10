package storage

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
)

// type: local
// path: /data/backups
type Local struct {
}

func (ctx *Local) perform(model config.ModelConfig, fileKey, archivePath string) error {
	logger.Info("=> storage | Local")
	destPath := model.StoreWith.Viper.GetString("path")
	helper.MkdirP(destPath)
	_, err := helper.Exec("cp", archivePath, destPath)
	if err != nil {
		return err
	}
	logger.Info("Store success", destPath)
	return nil
}
