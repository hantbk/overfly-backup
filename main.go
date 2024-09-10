package main

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/logger"
	"github.com/hantbk/vts-backup/model"
	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
	"os"
)

const (
	usage = "Backup Agent"
)

var (
	modelName  = ""
	configFile = ""
	version    = "master"
)

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = "vts-backup"
	app.Usage = usage

	app.Commands = []cli.Command{
		{
			Name: "perform",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "model, m",
					Usage:       "Model name that you want execute",
					Destination: &modelName,
				},
				cli.StringFlag{
					Name:        "config, c",
					Usage:       "Special a config file",
					Destination: &configFile,
				},
			},
			Action: func(c *cli.Context) error {
				config.Init(configFile)
				if len(modelName) == 0 {
					performAll()
				} else {
					performOne(modelName)
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func performAll() {
	for _, modelConfig := range config.Models {
		m := model.Model{
			Config: modelConfig,
		}
		m.Perform()
	}
}

func performOne(modelName string) {
	for _, modelConfig := range config.Models {
		if modelConfig.Name == modelName {
			m := model.Model{
				Config: modelConfig,
			}
			m.Perform()
			return
		}
	}
	logger.Fatalf("Model %s not found in %s", modelName, viper.ConfigFileUsed())
}
