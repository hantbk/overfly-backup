package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/logger"
	"github.com/hantbk/vtsbackup/model"
	"github.com/hantbk/vtsbackup/scheduler"
	"github.com/hantbk/vtsbackup/web"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	usage = "Backup agent."
)

var (
	configFile string
	version    = "master"
	signal     = flag.String("s", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown
  reload — reloading the configuration file`)
)

func buildFlags(flags []cli.Flag) []cli.Flag {
	return append(flags, &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Special a config file",
		Destination: &configFile,
	})
}

func termHandler(sig os.Signal) error {
	logger.Info("Received QUIT signal, exiting...")
	scheduler.Stop()
	os.Exit(0)
	return nil
}

func reloadHandler(sig os.Signal) error {
	logger.Info("Reloading config...")
	err := config.Init(configFile)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func main() {
	app := cli.NewApp()

	app.Version = version
	app.Name = "vtsbackup"
	app.Usage = usage

	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	app.Commands = []*cli.Command{
		{
			Name: "perform",
			Flags: buildFlags([]cli.Flag{
				&cli.StringSliceFlag{
					Name:    "model",
					Aliases: []string{"m"},
					Usage:   "Model name that you want perform",
				},
			}),
			Action: func(ctx *cli.Context) error {
				var modelNames []string
				err := initApplication()
				if err != nil {
					return err
				}
				modelNames = append(ctx.StringSlice("model"), ctx.Args().Slice()...)

				return perform(modelNames)
			},
		},
		{
			Name:  "start",
			Usage: "Start as daemon",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				fmt.Println("VtsBackup starting...")

				args := []string{"vtsbackup", "run"}
				if len(configFile) != 0 {
					args = append(args, "--config", configFile)
				}

				dm := &daemon.Context{
					// LogFileName: config.LogFilePath,
					PidFileName: config.PidFilePath,
					PidFilePerm: 0644,
					WorkDir:     "./",
					Args:        args,
				}

				d, err := dm.Reborn()
				if err != nil {
					return fmt.Errorf("start failed, please check is there another instance running: %w", err)
				}
				if d != nil {
					return nil
				}
				defer dm.Release() //nolint: errcheck

				logger.SetLogger(config.LogFilePath)

				err = initApplication()
				if err != nil {
					return err
				}

				if err := scheduler.Start(); err != nil {
					return fmt.Errorf("scheduler start failed: %w", err)
				}

				return nil
			},
		},
		{
			Name:  "run",
			Usage: "Run VtsBackup",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {

				logger.SetLogger(config.LogFilePath)

				err := initApplication()
				if err != nil {
					return err
				}

				if err := scheduler.Start(); err != nil {
					return fmt.Errorf("scheduler start failed: %w", err)
				}

				return web.StartHTTP(version)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}

func initApplication() error {
	return config.Init(configFile)
}

func perform(modelNames []string) error {
	var models []*model.Model
	if len(modelNames) == 0 {
		// perform all
		models = model.GetModels()
	} else {
		for _, name := range modelNames {
			if m := model.GetModelByName(name); m == nil {
				return fmt.Errorf("model %s not found in %s", name, viper.ConfigFileUsed())
			} else {
				models = append(models, m)
			}
		}
	}

	for _, m := range models {
		if err := m.Perform(); err != nil {
			logger.Tag(fmt.Sprintf("Model %s", m.Config.Name)).Error(err)
		}
	}
	return nil
}
