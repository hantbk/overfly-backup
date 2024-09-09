package main

import (
	"github.com/hantbk/vts-backup/config"
	"gopkg.in/urfave/cli.v1"
	"os"
)

const (
	usage = "Backup operation"
)

var (
	modelName = "all"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "vts-backup"
	app.Usage = usage

	app.Commands = []cli.Command{
		cli.Command{
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
					modelName = "all"
				}

				if modelName == "all" {
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
		model := Model{
			Config: modelConfig,
		}
		model.perform()
	}
}

func performOne(modelName string) {
	for _, modelConfig := range config.Models {
		if modelConfig.Name == modelName {
			model := Model{
				Config: modelConfig,
			}
			model.perform()
			return
		}
	}
}
