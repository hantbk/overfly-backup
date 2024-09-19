package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
			Name:  "backup",
			Usage: "Perform a full server backup",
			Flags: buildFlags([]cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "Output directory for the backup (default: ~/backups)",
					Value:   filepath.Join(os.Getenv("HOME"), "backups"),
				},
			}),
			Action: func(ctx *cli.Context) error {
				outputDir := ctx.String("output")
				return performFullBackup(outputDir)
			},
		},
		{
			Name:  "restore",
			Usage: "Restore server using config file",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				if len(configFile) == 0 {
					return fmt.Errorf("config file is required")
				}
				return restoreServer(configFile)
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

// func getDisks() ([]string, error) {
// 	var cmd *exec.Cmd
// 	osinfo, _, err := helper.CheckOS()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if osinfo == "darwin" {
// 		cmd = exec.Command("sh", "-c", "diskutil list | grep '/dev/disk' | awk '{print $1}'")
// 	} else if osinfo == "linux" {
// 		cmd = exec.Command("sh", "-c", "lsblk -d -o NAME | grep -E '^(nvme|sd|vd|xvd|hd)'")
// 	} else {
// 		return nil, fmt.Errorf("unsupported OS: %s", osinfo)
// 	}

// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	disks := strings.Fields(string(output))
// 	return disks, nil
// }

func performFullBackup(outputDir string) error {
	fmt.Printf("Performing full server backup to %s\n", outputDir)

	// Create necessary directories
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Path to the backup script
	backupScript := "./backup.sh"

	// Ensure the script is executable
	err = os.Chmod(backupScript, 0755)
	if err != nil {
		return fmt.Errorf("failed to make backup script executable: %w", err)
	}

	// Execute the backup script with the configuration file
	cmd := exec.Command(backupScript, configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute backup script: %w", err)
	}

	return nil
}

func restoreServer(configFile string) error {
	fmt.Printf("Restoring server using config file %s\n", configFile)

	// Path to the restore script
	restoreScript := "./restore.sh"

	// Ensure the script is executable
	err := os.Chmod(restoreScript, 0755)
	if err != nil {
		return fmt.Errorf("failed to make restore script executable: %w", err)
	}

	// Execute the restore script with the configuration file
	cmd := exec.Command(restoreScript, configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute restore script: %w", err)
	}

	return nil
}

// func performFullBackup(outputDir string) error {
// 	fmt.Printf("Performing full server backup to %s\n", outputDir)

// 	// Create necessary directories
// 	err := os.MkdirAll(outputDir, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	// Backup files and folders
// 	err = backupFiles(outputDir)
// 	if err != nil {
// 		return fmt.Errorf("failed to backup files: %w", err)
// 	}

// 	// Backup server configuration
// 	err = backupConfig(outputDir)
// 	if err != nil {
// 		return fmt.Errorf("failed to backup config: %w", err)
// 	}

// 	// Backup websites
// 	err = backupWebsites(outputDir)
// 	if err != nil {
// 		return fmt.Errorf("failed to backup websites: %w", err)
// 	}

// 	// Change permissions if required
// 	if config.MaintainFilePermissions == 0 {
// 		err = changePermissions(outputDir)
// 		if err != nil {
// 			return fmt.Errorf("failed to change permissions: %w", err)
// 		}
// 	}

// 	return nil
// }

// func backupFiles(outputDir string) error {
// 	fmt.Println("Backing up files and folders...")
// 	for _, resource := range config.BackupList {
// 		fmt.Printf("Creating backup for '%s'...\n", resource)
// 		if _, err := os.Stat(resource); os.IsNotExist(err) {
// 			fmt.Printf("Resource '%s' doesn't exist\n", resource)
// 			continue
// 		}
// 		cmd := exec.Command("rsync", "-avz", "--delete", "--relative", "--exclude-from", createExcludeFile(), resource, filepath.Join(outputDir, "files"))
// 		err := cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("rsync error: %w", err)
// 		}
// 	}
// 	return nil
// }

// func createExcludeFile() string {
// 	excludeFile := "/tmp/_VTS_backup_excluded"
// 	file, err := os.Create(excludeFile)
// 	if err != nil {
// 		fmt.Println("Error creating exclude file:", err)
// 		return ""
// 	}
// 	defer file.Close()

// 	for _, exclude := range config.BackupListExcluded {
// 		file.WriteString(exclude + "\n")
// 	}
// 	file.WriteString(config.BackupDir + "\n")
// 	return excludeFile
// }

// func backupConfig(outputDir string) error {
// 	fmt.Println("Backing up server configuration...")
// 	configDir := filepath.Join(outputDir, "config")
// 	err := os.MkdirAll(configDir, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create config directory: %w", err)
// 	}

// 	if config.CheckUsers == 1 {
// 		fmt.Println("Creating backup for users...")
// 		userBackupFile := filepath.Join(configDir, "users")
// 		cmd := exec.Command("cp", "/etc/passwd", userBackupFile)
// 		err := cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error backing up users: %w", err)
// 		}
// 	}

// 	if config.BackupIptables == 1 {
// 		fmt.Println("Creating backup for iptables...")
// 		cmd := exec.Command("/sbin/iptables-save", ">", filepath.Join(configDir, "iptables"))
// 		err := cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error backing up iptables: %w", err)
// 		}
// 	}

// 	if config.BackupCrontabs == 1 {
// 		fmt.Println("Creating backup for crontabs...")
// 		crontabsDir := filepath.Join(configDir, "crontabs")
// 		err := os.MkdirAll(crontabsDir, 0755)
// 		if err != nil {
// 			return fmt.Errorf("failed to create crontabs directory: %w", err)
// 		}
// 		cmd := exec.Command("cp", "-f", "/var/spool/cron/*", crontabsDir)
// 		err = cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error backing up crontabs: %w", err)
// 		}
// 	}

// 	if config.CheckEnabledServices == 1 {
// 		fmt.Println("Checking enabled services...")
// 		cmd := exec.Command("chkconfig", "--list", "|", "grep", "-i", "2:activ", "|", "awk", "'{print $1}'", "|", "sort", ">", filepath.Join(configDir, "enabledServices"))
// 		err := cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error checking enabled services: %w", err)
// 		}
// 	}

// 	if config.CheckInstalledPackages == 1 {
// 		fmt.Println("Checking installed packages...")
// 		cmd := exec.Command("rpm", "-qa", "--qf", "%{NAME}\n", "|", "sort", ">", "/tmp/installedPackages.tmp")
// 		err := cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error checking installed packages: %w", err)
// 		}
// 		cmd = exec.Command("sh", "-c", "while read concretePackage; do provider=`rpm -q --whatprovides \"$concretePackage\" --qf \"%{NAME}\n\" | head -n 1`; echo \"$concretePackage\t$provider\" >> "+filepath.Join(configDir, "installedPackages")+"; done < /tmp/installedPackages.tmp")
// 		err = cmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error adding package information: %w", err)
// 		}
// 	}

// 	return nil
// }

// func backupWebsites(outputDir string) error {
// 	// Implement the logic to backup websites
// 	fmt.Println("Backing up websites...")
// 	// Example: rsync command to backup websites
// 	cmd := exec.Command("rsync", "-avz", "--delete", "/var/www", filepath.Join(outputDir, "websites"))
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("rsync error: %w", err)
// 	}
// 	return nil
// }

// func changePermissions(outputDir string) error {
// 	fmt.Println("Changing permissions to allow external users to access the backups...")
// 	cmd := exec.Command("find", outputDir, "-type", "d", "-exec", "chmod", "o+xr", "{}", ";")
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("error changing directory permissions: %w", err)
// 	}
// 	cmd = exec.Command("find", outputDir, "-type", "f", "-exec", "chmod", "o+r", "{}", ";")
// 	err = cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("error changing file permissions: %w", err)
// 	}
// 	return nil
// }
