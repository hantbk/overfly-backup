package storage

import "github.com/hantbk/vts-backup/config"

// Base storage
type Base interface {
	perform(model config.ModelConfig, archivePath string) error
}

// Run storage
func Run(model config.ModelConfig, archivePath string) error {
	var ctx Base
	switch model.StoreWith.Type {
	case "local":
		ctx = &Local{}
	case "ftp":
		ctx = &FTP{}
	default:
		ctx = &Local{}
	}
	err := ctx.perform(model, archivePath)
	if err != nil {
		return err
	}
	return nil
}
