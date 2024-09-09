package compressor

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/logger"
)

type Base interface {
	perform(model config.ModelConfig) (resultPath *string, err error)
}

// Run compress
func Run(model config.ModelConfig) (resultPath *string, err error) {
	logger.Info("----------------Compressing----------------")
	var ctx Base
	switch model.CompressWith.Type {
	case "tgz":

		ctx = &Tgz{}
	default:
		ctx = &Tgz{}
	}

	resultPath, err = ctx.perform(model)
	if err != nil {
		return
	}
	logger.Info("----------------Compressing done----------------")

	return
}
