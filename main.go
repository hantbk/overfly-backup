package main

import (
	"github.com/hantbk/vts-backup/config"
	"github.com/hantbk/vts-backup/model"
	"gopkg.in/urfave/cli.v1"
	"os"
)

const (
	usage = "Backup Agent"
)

var (
	modelName = ""
	version   = "0.0.1"
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
			},
			Action: func(c *cli.Context) error {
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
}
