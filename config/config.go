package config

import (
	"fmt"
	"github.com/hantbk/vts-backup/logger"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

var (
	// Exist Is config file exist
	Exist bool
	// Models configs
	Models []ModelConfig
	// HomeDir of user
	HomeDir = os.Getenv("HOME")
)

// ModelConfig for special case
type ModelConfig struct {
	Name         string
	TempPath     string
	DumpPath     string
	CompressWith SubConfig
	EncryptWith  SubConfig
	StoreWith    SubConfig
	Archive      *viper.Viper
	Databases    []SubConfig
	Storages     []SubConfig
	Viper        *viper.Viper
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
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
		logger.Error("Load backup config fail", err)
		return
	}

	Exist = true
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}

}

func loadModel(key string) (model ModelConfig) {
	model.Name = key
	model.TempPath = path.Join(os.TempDir(), "vtsbackup", fmt.Sprintf("%d", time.Now().UnixNano()))
	model.DumpPath = path.Join(model.TempPath, key)
	model.Viper = viper.Sub("models." + key)

	model.CompressWith = SubConfig{
		Type:  model.Viper.GetString("compress_with.type"),
		Viper: model.Viper.Sub("compress_with"),
	}

	model.EncryptWith = SubConfig{
		Type:  model.Viper.GetString("encrypt_with.type"),
		Viper: model.Viper.Sub("encrypt_with"),
	}

	model.StoreWith = SubConfig{
		Type:  model.Viper.GetString("store_with.type"),
		Viper: model.Viper.Sub("store_with"),
	}

	model.Archive = model.Viper.Sub("archive")

	loadStoragesConfig(&model)

	return
}

func loadStoragesConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("storages")
	for key := range model.Viper.GetStringMap("storages") {
		dbViper := subViper.Sub(key)
		model.Storages = append(model.Storages, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
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
