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
)

// ModelConfig for special case
type ModelConfig struct {
	Name         string
	DumpPath     string
	CompressWith SubConfig
	EncryptWith  SubConfig
	StoreWith    SubConfig
	Archive      *viper.Viper
	Storages     []SubConfig
	Viper        *viper.Viper
}

// Subconfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	// ./vts-backup.yml
	viper.AddConfigPath(".")
	// ~/.vts-backup/vts-backup.yml
	viper.AddConfigPath("$HOME/.vts-backup") // call multiple times to add many search paths
	// /etc/vts-backup/vts-backup.yml
	viper.AddConfigPath("/etc/vts-backup/") // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Load vts-backup config faild", err)
		return
	}
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}
	return
}

func init() {
	viper.SetConfigType("yaml")

	isTest := os.Getenv("GO_ENV") == "test"

	if isTest {

		viper.SetConfigName("config_test")
	} else {
		viper.SetConfigName("config")
	}

	// ./vts-backup.yml
	viper.AddConfigPath(".")
	if isTest {
		viper.AddConfigPath("../")
	} else {
		// ~/.vts-backup/vts-backup.yml
		viper.AddConfigPath("$HOME/.vts-backup") // call multiple times to add many search paths
		// /etc/vts-backup/vts-backup.yml
		viper.AddConfigPath("/etc/vts-backup/") // path to look for the config file in
	}

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Load vts-backup config faild", err)
		return
	}
	Exist = true
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}
	return
}

func loadModel(key string) (model ModelConfig) {
	model.Name = key
	model.DumpPath = path.Join(os.TempDir(), "vts-backup", fmt.Sprintf("%d", time.Now().UnixNano()), key)
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
