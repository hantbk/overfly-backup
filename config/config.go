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
	DumpPath     string
	CompressWith string
	Storages     []SubConfig
)

type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	// /etc/vts-backup/config.yaml
	viper.AddConfigPath("/etc/vts-backup")
	// ~/.vts-backup/config.yaml
	viper.AddConfigPath("$HOME/.vts-backup")

	// ./config.yaml
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Failed to read config file: ", err)
		return
	}

	DumpPath = path.Join(os.TempDir(), "vts-backup", fmt.Sprintf("%d", time.Now().UnixNano()))
	CompressWith = viper.GetString("compress_with")
	loadStoragesConfig()

	return
}

func loadStoragesConfig() {
	subViper := viper.Sub("storages")
	for key := range viper.GetStringMap("storages") {
		dbViper := subViper.Sub(key)
		Storages = append(Storages, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
}
