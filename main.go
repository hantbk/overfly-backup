package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/helper"
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
			Name:  "perform",
			Usage: "Perform backup pipeline using config file",
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
			Usage: "Start Backup agent as daemon",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				fmt.Println("Backup starting as daemon...")

				args := []string{"vtsbackup", "run"}
				if len(configFile) != 0 {
					args = append(args, "--config", configFile)
				}

				dm := &daemon.Context{
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
				defer dm.Release() //nolint:errcheck

				logger.SetLogger(config.LogFilePath)

				err = initApplication()
				if err != nil {
					return err
				}

				if err := scheduler.Start(); err != nil {
					return fmt.Errorf("failed to start scheduler: %w", err)
				}

				return nil
			},
		},
		{
			Name:  "run",
			Usage: "Run Backup agent without daemon",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				logger.SetLogger(config.LogFilePath)

				err := initApplication()
				if err != nil {
					return err
				}

				if err := scheduler.Start(); err != nil {
					return fmt.Errorf("failed to start scheduler: %w", err)
				}

				return web.StartHTTP(version)
			},
		},
		{
			Name:  "list",
			Usage: "List running Backup agents",
			Action: func(ctx *cli.Context) error {
				pids, err := listBackupAgents()
				if err != nil {
					return err
				}
				if len(pids) == 0 {
					fmt.Println("No running Backup agents found.")
				} else {
					fmt.Println("Running Backup agents PIDs:")
					for _, pid := range pids {
						fmt.Println(pid)
					}
				}
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop the running Backup agent",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Stopping Backup agent...")
				pids, err := listBackupAgents()
				if err != nil {
					return err
				}
				if len(pids) == 0 {
					fmt.Println("No running Backup agents found.")
				} else {
					fmt.Println("Running Backup agents PIDs:")
					for _, pid := range pids {
						fmt.Println(pid)
					}
				}
				for _, pid := range pids {
					stopBackupAgent(pid)
				}
				return nil
			},
		},
		{
			Name:  "reload",
			Usage: "Reload the running Backup agent",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Reloading Backup agent...")
				pids, err := listBackupAgents()
				if err != nil {
					return err
				}

				if len(pids) == 0 {
					fmt.Println("No running Backup agents found.")
				} else {
					fmt.Println("Running Backup agents PIDs:")
					for _, pid := range pids {
						fmt.Println(pid)
					}
				}
				for _, pid := range pids {
					reloadBackupAgent(pid)
				}
				return nil
			},
		},
		{
			Name:  "listdisk",
			Usage: "List all disk on machine.",
			Action: func(ctx *cli.Context) error {
				// fmt.Println("All disks:")
				listDisk()
				return nil
			},
		},
		{
			Name:  "createimg",
			Usage: "Create disk image",
			Flags: buildFlags([]cli.Flag{
				&cli.StringFlag{
					Name:    "disk",
					Aliases: []string{"d"},
					Usage:   "Disk name to create image for",
				},
			}),
			Action: func(ctx *cli.Context) error {
				disk := ctx.String("disk")
				createDiskImage(disk)
				return nil
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

func listBackupAgents() ([]int, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	var pids []int
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "vtsbackup") && !strings.Contains(line, "grep") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid, err := strconv.Atoi(fields[1])
				if err == nil {
					pids = append(pids, pid)
				}
			}
		}
	}
	return pids, nil
}

func stopBackupAgent(pid int) {
	cmd := exec.Command("kill", "-QUIT", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Backup agent stopped successfully")
	}
}

func reloadBackupAgent(pid int) {
	cmd := exec.Command("kill", "-HUP", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Backup agent reloaded successfully")
	}
}

func listDisk() {
	osinfo, _, err := helper.CheckOS()
	if err != nil {
		return
	}
	var cmd *exec.Cmd
	if osinfo == "darwin" {
		cmd = exec.Command("diskutil", "list")
	} else if osinfo == "linux" {
		cmd = exec.Command("lsblk")
	} else {
		fmt.Println("Unsupported OS:", osinfo)
		return
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func getDisks() ([]string, error) {
	var cmd *exec.Cmd
	osinfo, _, err := helper.CheckOS()
	if err != nil {
		return nil, err
	}

	if osinfo == "darwin" {
		cmd = exec.Command("sh", "-c", "diskutil list | grep '/dev/disk' | awk '{print $1}'")
	} else if osinfo == "linux" {
		cmd = exec.Command("sh", "-c", "lsblk -d -o NAME | grep -E '^(nvme|sd|vd|xvd|hd)'")
	} else {
		return nil, fmt.Errorf("unsupported OS: %s", osinfo)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	disks := strings.Fields(string(output))
	fmt.Println(disks)
	return disks, nil
}

func createDiskImage(option string) {
	disks, err := getDisks()
	if err != nil {
		fmt.Println("Error getting disks:", err)
		return
	}

	backupDir := os.ExpandEnv("$HOME/backups")
	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		fmt.Println("Error creating backup directory:", err)
		return
	}

	if option == "" {
		for _, disk := range disks {
			imagePath := fmt.Sprintf("%s/%s_backup.img", backupDir, disk)
			cmd := exec.Command("sudo", "dd", "if=/dev/"+disk, "of="+imagePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Println("Error creating disk image for", disk, ":", err)
			} else {
				fmt.Println("Disk image created successfully:", imagePath)
			}
		}
	} else {
		if !contains(disks, option) {
			fmt.Printf("Disk %s not found.\n", option)
			return
		}

		imagePath := fmt.Sprintf("%s/%s_backup.img", backupDir, option)
		cmd := exec.Command("sudo", "dd", "if=/dev/"+option, "of="+imagePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("Error creating disk image for", option, ":", err)
		} else {
			fmt.Println("Disk image created successfully:", imagePath)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
