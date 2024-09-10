package model

import (
	"fmt"
	"github.com/hantbk/vts-backup/archive"
	"github.com/hantbk/vts-backup/compressor"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/encryptor"
	"github.com/hantbk/vts-backup/logger"
	"github.com/hantbk/vts-backup/splitter"
	"github.com/hantbk/vts-backup/storage"
	"github.com/spf13/viper"
	"os"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (m Model) Perform() {
	logger := logger.Tag(fmt.Sprintf("Model: %s", m.Config.Name))
	logger.Info("WorkDir:", m.Config.DumpPath)

	defer func() {
		if r := recover(); r != nil {
			m.cleanup()
		}

		m.cleanup()
	}()

	if m.Config.Archive != nil {
		err := archive.Run(m.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	archivePath, err := compressor.Run(m.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = encryptor.Run(archivePath, m.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = splitter.Run(archivePath, m.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(m.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

}

// Cleanup model temp files
func (m Model) cleanup() {
	logger := logger.Tag("Modal")
	tempDir := m.Config.TempPath
	if viper.GetBool("useTempWorkDir") {
		tempDir = viper.GetString("workdir")
	}
	logger.Infof("Cleanup temp: %s/", tempDir)
	if err := os.RemoveAll(tempDir); err != nil {
		logger.Errorf("Cleanup temp dir %s error: %v", tempDir, err)
	}
}
