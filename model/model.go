package model

import (
	"fmt"
	"github.com/hantbk/vts-backup/archive"
	"github.com/hantbk/vts-backup/compressor"
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/encryptor"
	"github.com/hantbk/vts-backup/logger"
	"github.com/hantbk/vts-backup/notifier"
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
func (m Model) Perform() (err error) {
	logger := logger.Tag(fmt.Sprintf("Model: %s", m.Config.Name))

	defer func() {
		if err != nil {
			logger.Error(err)
			notifier.Failure(m.Config, err.Error())
		} else {
			notifier.Success(m.Config)
		}
	}()

	logger.Info("WorkDir:", m.Config.DumpPath)

	defer func() {
		if r := recover(); r != nil {
			m.cleanup()
		}

		m.cleanup()
	}()

	if m.Config.Archive != nil {
		err = archive.Run(m.Config)
		if err != nil {
			return
		}
	}

	archivePath, err := compressor.Run(m.Config)
	if err != nil {
		return
	}

	archivePath, err = encryptor.Run(archivePath, m.Config)
	if err != nil {
		return
	}

	archivePath, err = splitter.Run(archivePath, m.Config)
	if err != nil {
		return
	}

	err = storage.Run(m.Config, archivePath)
	if err != nil {
		return
	}

	return nil
}

// Cleanup model temp files
func (m Model) cleanup() {
	logger := logger.Tag("Model")

	tempDir := m.Config.TempPath
	if viper.GetBool("useTempWorkDir") {
		tempDir = viper.GetString("workdir")
	}
	logger.Infof("Cleanup temp: %s/", tempDir)
	if err := os.RemoveAll(tempDir); err != nil {
		logger.Errorf("Cleanup temp dir %s error: %v", tempDir, err)
	}
}

// GetModelByName get model by name
func GetModelByName(name string) *Model {
	modelConfig := config.GetModelConfigByName(name)
	if modelConfig == nil {
		return nil
	}
	return &Model{
		Config: *modelConfig,
	}
}

// GetModels get models
func GetModels() (models []*Model) {
	for _, modelConfig := range config.Models {
		m := Model{
			Config: modelConfig,
		}
		models = append(models, &m)
	}
	return
}
