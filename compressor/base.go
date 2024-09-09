package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/logger"
)

type Base interface {
	perform() error
}

// Run compress
func Run() error {
	logger.Info("----------------Compressing----------------")
	var ctx Base
	switch config.CompressWith {
	case "tgz":
		ctx = &Tgz{}
	default:
		ctx = &Tgz{}
	}

	err := ctx.perform()
	if err != nil {
		return err
	}
	logger.Info("----------------Compressing done----------------")

	return nil
}
