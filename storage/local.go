package storage

import "github.com/hantbk/vts-backup/logger"

type Local struct {
}

func newLocal() *Local {
	return &Local{}
}

func (ctx *Local) Perform() error {
	logger.Info("Performing local storage")
	return nil
}
