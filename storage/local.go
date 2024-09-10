package storage

import (
	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"
	"path"
)

// Local storage
//
// type: local
// path: /data/backups
type Local struct {
	Base
	destPath string
}

func (s *Local) open() error {
	s.destPath = s.viper.GetString("path")
	return helper.MkdirP(s.destPath)
}

func (s *Local) close() {}

func (s *Local) upload(fileKey string) (err error) {
	logger := logger.Tag("Local")

	_, err = helper.Exec("cp", s.archivePath, s.destPath)
	if err != nil {
		return err
	}
	logger.Info("Store successed", s.destPath)
	return nil
}

func (s *Local) delete(fileKey string) (err error) {
	_, err = helper.Exec("rm", path.Join(s.destPath, fileKey))
	return
}
