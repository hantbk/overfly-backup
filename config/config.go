package config

import (
	"fmt"
	"github.com/hantbk/vts-backup/logger"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

var (
	// Exist Is config file exist
	Exist bool
	// Models configs
	Models []ModelConfig
	// vtsbackup base dir
	VtsBackupDir string = getVtsBackupDir()
)

// ModelConfig for special case
type ModelConfig struct {
	Name         string
	TempPath     string
	DumpPath     string
	CompressWith SubConfig
	EncryptWith  SubConfig
	Archive      *viper.Viper
	Splitter     *viper.Viper
	Storages     map[string]SubConfig
	Viper        *viper.Viper
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func getVtsBackupDir() string {
	dir := os.Getenv("VTSBACKUP_DIR")
	if len(dir) == 0 {
		dir = filepath.Join(os.Getenv("HOME"), ".vtsbackup")
	}
	return dir
}

// loadConfig from:
// - ./vtsbackup.yml
// - ~/.vtsbackup/vtsbackup.yml
// - /etc/vtsbackup/vtsbackup.yml
func Init(configFile string) {
	viper.SetConfigType("yaml")

	// Set config file directly
	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("vtsbackup") // name of config file (without extension)

		// ./vtsbackup.yml
		viper.AddConfigPath(".")

		// ~/.vtsbackup/vtsbackup.yml
		viper.AddConfigPath("$HOME/.vtsbackup")

		// /etc/vtsbackup/vtsbackup.yml
		viper.AddConfigPath("/etc/vtsbackup")

	}

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Load backup config fail: ", err)
		return
	}

	viperConfigFile := viper.ConfigFileUsed()
	if info, _ := os.Stat(viperConfigFile); info.Mode()&(1<<2) != 0 {
		// max permission: 0770
		logger.Warnf("Other users are able to access %s with mode %v", viperConfigFile, info.Mode())
	}

	viper.Set("useTempWorkDir", false)
	if workdir := viper.GetString("workdir"); len(workdir) == 0 {
		// use temp dir as workdir
		dir, err := os.MkdirTemp("", "vtsbackup")
		if err != nil {
			logger.Fatal(err)
		}
		viper.Set("workdir", dir)
		viper.Set("useTempWorkDir", true)
	}

	Exist = true
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}

	if len(Models) == 0 {
		logger.Fatalf("No model found in %s", viperConfigFile)
	}
}

func loadModel(key string) (model ModelConfig) {
	model.Name = key
	model.TempPath = filepath.Join(viper.GetString("workdir"), fmt.Sprintf("%d", time.Now().UnixNano()))
	model.DumpPath = filepath.Join(model.TempPath, key)
	model.Viper = viper.Sub("models." + key)

	model.CompressWith = SubConfig{
		Type:  model.Viper.GetString("compress_with.type"),
		Viper: model.Viper.Sub("compress_with"),
	}

	model.EncryptWith = SubConfig{
		Type:  model.Viper.GetString("encrypt_with.type"),
		Viper: model.Viper.Sub("encrypt_with"),
	}

	model.Archive = model.Viper.Sub("archive")

	model.Splitter = model.Viper.Sub("split_with")

	loadStoragesConfig(&model)

	if len(model.Storages) == 0 {
		logger.Fatalf("No storage found in model %s", model.Name)
	}

	return
}

func loadStoragesConfig(model *ModelConfig) {
	storageConfigs := map[string]SubConfig{}
	// Backward compatible with `store_with` config
	storeWith := model.Viper.Sub("store_with")

	if storeWith != nil {
		logger.Warn(`[Deprecated] "store_with" is deprecated now, please use "storages" which supports multiple storages.`)
		storageConfigs["store_with"] = SubConfig{
			Name:  "",
			Type:  model.Viper.GetString("store_with.type"),
			Viper: model.Viper.Sub("store_with"),
		}
	}

	subViper := model.Viper.Sub("storages")

	for key := range model.Viper.GetStringMap("storages") {
		storageViper := subViper.Sub(key)
		storageConfigs[key] = SubConfig{
			Name:  key,
			Type:  storageViper.GetString("type"),
			Viper: storageViper,
		}
	}
	model.Storages = storageConfigs
}

// GetModelByName get model by name
func GetModelByName(name string) (model *ModelConfig) {
	for _, m := range Models {
		if m.Name == name {
			model = &m
			return
		}
	}
	return
}
