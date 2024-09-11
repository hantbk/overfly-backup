package storage

import (
	"fmt"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/logger"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Base storage
// When `archivePath` is a directory, `fileKeys` stores files in the `archivePath` with directory prefix
type Base struct {
	model       config.ModelConfig
	archivePath string
	fileKeys    []string
	viper       *viper.Viper
	keep        int
	cycler      *Cycler
}

// Storage interface
type Storage interface {
	open() error
	close()
	upload(fileKey string) error
	delete(fileKey string) error
}

func newBase(model config.ModelConfig, archivePath string, storageConfig config.SubConfig) (base Base, err error) {
	// Backward compatible with `store_with` config
	var cyclerName string
	if storageConfig.Name == "" {
		cyclerName = model.Name
	} else {
		cyclerName = fmt.Sprintf("%s_%s", model.Name, storageConfig.Name)
	}

	var keys []string
	if fi, err := os.Stat(archivePath); err == nil && fi.IsDir() {
		// NOTE: ignore err is not nil scenario here to pass test and should be fine
		// 2022.12.04.07.09.47
		entries, err := os.ReadDir(archivePath)
		if err != nil {
			return base, err
		}
		for _, e := range entries {
			// Assume all entries are file
			// 2022.12.04.07.09.47/2022.12.04.07.09.47.tar.xz-000
			if !e.IsDir() {
				keys = append(keys, filepath.Join(filepath.Base(archivePath), e.Name()))
			}
		}
	}

	base = Base{
		model:       model,
		archivePath: archivePath,
		fileKeys:    keys,
		viper:       storageConfig.Viper,
		cycler:      &Cycler{name: cyclerName},
	}

	if base.viper != nil {
		base.keep = base.viper.GetInt("keep")
	}

	return
}

// run storage
func runModel(model config.ModelConfig, archivePath string, storageConfig config.SubConfig) (err error) {
	logger := logger.Tag("Storage")

	newFileKey := filepath.Base(archivePath)
	base, err := newBase(model, archivePath, storageConfig)
	if err != nil {
		return err
	}
	var s Storage
	switch storageConfig.Type {
	case "local":
		s = &Local{Base: base}
	case "webdav":
		s = &WebDAV{Base: base}
	case "ftp":
		s = &FTP{Base: base}
	case "scp":
		s = &SCP{Base: base}
	case "s3":
		s = &S3{Base: base, Service: "s3"}
	default:
		return fmt.Errorf("[%s] storage type has not implement", storageConfig.Type)
	}

	logger.Info("=> Storage | " + storageConfig.Type)
	err = s.open()
	if err != nil {
		return err
	}
	defer s.close()

	err = s.upload(newFileKey)
	if err != nil {
		return err
	}

	base.cycler.run(newFileKey, base.fileKeys, base.keep, s.delete)

	return nil
}

// Run storage
func Run(model config.ModelConfig, archivePath string) (err error) {
	for _, storageConfig := range model.Storages {
		err := runModel(model, archivePath, storageConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
