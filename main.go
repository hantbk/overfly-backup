// Copyright © 2024 Ha Nguyen <captainnemot1k60@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/decompressor"
	"github.com/hantbk/vtsbackup/helper"
	"github.com/hantbk/vtsbackup/logger"
	"github.com/hantbk/vtsbackup/model"
	"github.com/hantbk/vtsbackup/scheduler"
	"github.com/hantbk/vtsbackup/storage"
	"github.com/hantbk/vtsbackup/web"
	"github.com/schollz/progressbar/v3"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

//go:embed vts
var vtsScript embed.FS

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
			Usage: "Perform backup pipeline using config file. If no model is specified, all models will be performed.",
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
		// {
		// 	Name:  "list-agent",
		// 	Usage: "List running Backup agents",
		// 	Action: func(ctx *cli.Context) error {
		// 		pids, err := listBackupAgents()
		// 		if err != nil {
		// 			return err
		// 		}
		// 		if len(pids) == 0 {
		// 			fmt.Println("No running Backup agents found.")
		// 		} else {
		// 			fmt.Println("Running Backup agents PIDs:")
		// 			for _, pid := range pids {
		// 				fmt.Println(pid)
		// 			}
		// 		}
		// 		return nil
		// 	},
		// },
		{
			Name:  "stop",
			Usage: "Stop the running Backup agent",
			Action: func(c *cli.Context) error {
				// fmt.Println("Stopping Backup agent...")
				err := stopBackupAgent()
				if err != nil {
					return cli.Exit(err.Error(), 1)
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
			Name:  "list",
			Usage: "List current saved backups",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				return runBashCommand("list")
			},
		},
		{
			Name:  "listM",
			Usage: "List all configured backup models",
			Flags: buildFlags([]cli.Flag{}),
			Action: func(ctx *cli.Context) error {
				return listModel()
			},
		},
		{
			Name:  "listB",
			Usage: "List backup files for a specific model in S3",
			Flags: buildFlags([]cli.Flag{
				&cli.StringFlag{
					Name:     "model",
					Aliases:  []string{"m"},
					Usage:    "Model name to list backups for",
					Required: true,
				},
			}),
			Action: func(ctx *cli.Context) error {
				modelName := ctx.String("model")
				return listBackupFiles(modelName)
			},
		},
		{
			Name:  "download",
			Usage: "Download a backup file for a specific model",
			Flags: buildFlags([]cli.Flag{
				&cli.StringFlag{
					Name:     "model",
					Aliases:  []string{"m"},
					Usage:    "Model name to download backup from",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "Path to save backup file",
				},
			}),
			Action: func(ctx *cli.Context) error {
				modelName := ctx.String("model")
				outputPath := ctx.String("output")
				return downloadBackupFile(modelName, outputPath)
			},
		},
		{
			Name:  "save",
			Usage: "Perform standard rotating backups",
			Action: func(ctx *cli.Context) error {
				err := runBashCommand("save")
				if err != nil {
					return err
				}

				// // Initialize the application before performing
				// if err := initApplication(); err != nil {
				// 	return err
				// }

				// // Execute perform command with sudo
				// cmd := exec.Command("sudo", os.Args[0], "perform")
				// cmd.Stdout = os.Stdout
				// cmd.Stderr = os.Stderr
				// return cmd.Run()
				return nil
			},
		},
		{
			Name:  "keep",
			Usage: "Perform manual backups (never deleted)",
			Action: func(ctx *cli.Context) error {
				err := runBashCommand("keep")
				if err != nil {
					return err
				}

				// // Initialize the application before performing
				// if err := initApplication(); err != nil {
				// 	return err
				// }

				// // Execute perform command with sudo
				// cmd := exec.Command("sudo", os.Args[0], "perform")
				// cmd.Stdout = os.Stdout
				// cmd.Stderr = os.Stderr
				// return cmd.Run()
				return nil
			},
		},
		{
			Name:  "mydata",
			Usage: "Backup user (personal) data - IN DEVELOPMENT",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("mydata", ctx.Args().First())
			},
		},
		{
			Name:  "restore",
			Usage: "Restore the system",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("restore")
			},
		},
		{
			Name:  "archive",
			Usage: "Tar and gzip profiles",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("archive")
			},
		},
		{
			Name:  "configfile",
			Usage: "Display configuration file help",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("configfile")
			},
		},
		{
			Name:  "directories",
			Usage: "List directories backed up",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("directories")
			},
		},
		{
			Name:  "script",
			Usage: "Automated Linux mounting script",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("script")
			},
		},
		{
			Name:  "update",
			Usage: "Display how to find updates",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("update")
			},
		},
		{
			Name:  "uses",
			Usage: "Some uses for this script",
			Action: func(ctx *cli.Context) error {
				return runBashCommand("uses")
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall backup agent",
			Action: func(ctx *cli.Context) error {
				return uninstallBackupAgent()
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

func stopBackupAgent() error {
	pids, err := findBackupAgentPIDs()
	if err != nil {
		return fmt.Errorf("error finding backup agent processes: %w", err)
	}

	// if len(pids) == 0 {
	// 	fmt.Println("No running backup agent processes found.")
	// 	return nil
	// }

	// fmt.Printf("Found %d running backup agent process(es).\n", len(pids))

	for _, pid := range pids {
		// fmt.Printf("Stopping process with PID %d...\n", pid)
		err := syscall.Kill(pid, syscall.SIGQUIT)
		if err != nil {
			return fmt.Errorf("error stopping process %d: %w", pid, err)
		}
	}

	fmt.Println("All backup agent processes have been stopped.")
	return nil
}

func findBackupAgentPIDs() ([]int, error) {
	cmd := exec.Command("pgrep", "-f", "vtsbackup.*run")
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			// No processes found, not an error
			return nil, nil
		}
		return nil, err
	}

	var pids []int
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		pid, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing PID: %w", err)
		}
		pids = append(pids, pid)
	}

	return pids, nil
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

func listModel() error {
	err := initApplication()
	if err != nil {
		return err
	}
	models := model.GetModels()
	if len(models) == 0 {
		fmt.Println("No backup models configured.")
	} else {
		fmt.Printf("Configured backup models (from %s):\n", viper.ConfigFileUsed())
		for _, m := range models {
			fmt.Printf("- %s\n", m.Config.Name)
			if m.Config.Description != "" {
				fmt.Printf("  Description: %s\n", m.Config.Description)
			}
			if m.Config.Schedule.Enabled {
				fmt.Printf("  Schedule: %s\n", m.Config.Schedule.String())
			}
			if m.Config.Archive != nil {
				fmt.Println("  Archive:")
				if includes := m.Config.Archive.GetStringSlice("includes"); len(includes) > 0 {
					fmt.Println("    Includes:")
					for _, include := range includes {
						fmt.Printf("      - %s\n", include)
					}
				}
				if excludes := m.Config.Archive.GetStringSlice("excludes"); len(excludes) > 0 {
					fmt.Println("    Excludes:")
					for _, exclude := range excludes {
						fmt.Printf("      - %s\n", exclude)
					}
				}
			}
			if m.Config.CompressWith != (config.SubConfig{}) {
				fmt.Printf("  Compression: %s\n", m.Config.CompressWith.Type)
			}
			if m.Config.EncryptWith != (config.SubConfig{}) {
				fmt.Printf("  Encryption: %s\n", m.Config.EncryptWith.Type)
			}
			if m.Config.Storages != nil {
				fmt.Println("  Storages:")
				for name, storage := range m.Config.Storages {
					fmt.Printf("    - %s:%s\n", name, storage.Type)
					if storage.Type == "local" {
						fmt.Printf("      Path: %s\n", storage.Viper.GetString("path"))
					} else if storage.Type == "s3" || storage.Type == "minio" {
						fmt.Printf("      Bucket: %s\n", storage.Viper.GetString("bucket"))
						fmt.Printf("      Path: %s\n", storage.Viper.GetString("path"))
					} else if storage.Type == "scp" {
						fmt.Printf("      Host: %s\n", storage.Viper.GetString("host"))
						fmt.Printf("      Path: %s\n", storage.Viper.GetString("path"))
					}
				}
			}
			fmt.Println()
		}
	}
	return nil
}

func listBackupFiles(modelName string) error {
	err := initApplication()
	if err != nil {
		return err
	}

	m := model.GetModelByName(modelName)
	if m == nil {
		return fmt.Errorf("model: %q not found", modelName)
	}

	files, err := storage.List(m.Config, "/")
	if err != nil {
		return fmt.Errorf("failed to list backup files: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("No backup files found for model %q\n", modelName)
	} else {
		fmt.Printf("Backup files for model %q:\n", modelName)
		for _, file := range files {
			fmt.Printf("- %s (Size: %s, Last Modified: %s)\n",
				file.Filename,
				humanize.Bytes(uint64(file.Size)),
				file.LastModified.Format(time.RFC3339),
			)
		}
	}

	return nil
}

func downloadBackupFile(modelName, outputPath string) error {
	if outputPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %v", err)
		}
		outputPath = homeDir
	} else {
		outputPath = helper.ExplandHome(outputPath)
	}
	err := initApplication()
	if err != nil {
		return err
	}

	m := model.GetModelByName(modelName)
	if m == nil {
		return fmt.Errorf("model: %q not found", modelName)
	}

	files, err := storage.List(m.Config, "/")
	if err != nil {
		return fmt.Errorf("failed to list backup files: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("No backup files found for model %q\n", modelName)
		return nil
	}

	fmt.Printf("Backup files for model %q:\n", modelName)
	for i, file := range files {
		fmt.Printf("%d. %s (Size: %s, Last Modified: %s)\n",
			i+1,
			file.Filename,
			humanize.Bytes(uint64(file.Size)),
			file.LastModified.Format(time.RFC3339),
		)
	}

	var choice int
	fmt.Print("Enter the number of the file you want to download (0 to cancel): ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 0 || choice > len(files) {
		return fmt.Errorf("invalid choice")
	}

	if choice == 0 {
		fmt.Println("Download cancelled.")
		return nil
	}

	selectedFile := files[choice-1]

	fmt.Printf("You selected: %s\n", selectedFile.Filename)
	fmt.Print("Do you want to proceed with the download? (Y/n): ")
	var confirm string
	fmt.Scanf("%s", &confirm)

	if strings.ToLower(confirm) != "y" && confirm != "" {
		fmt.Println("Download cancelled.")
		return nil
	}

	downloadURL, err := storage.Download(m.Config, selectedFile.Filename)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %v", err)
	}

	filePath := filepath.Join(outputPath, selectedFile.Filename)
	dirPath := filepath.Dir(filePath)

	// Ensure the output directory exists
	if err := helper.MkdirP(dirPath); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	fmt.Printf("Downloading %s to %s...\n", selectedFile.Filename, filePath)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Create a progress bar
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	// Write the body to file
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("\nFile downloaded successfully: %s\n", filePath)
	// Decompress file
	if err := decompressor.Run(filePath, modelName); err != nil {
		return fmt.Errorf("failed to decompress file: %v", err)
	}

	return nil
}

func uninstallBackupAgent() error {
	// fmt.Println("Uninstalling backup agent...")

	// Stop the daemon
	if err := stopBackupAgent(); err != nil {
		fmt.Printf("Warning: Failed to stop backup agent: %v\n", err)
	}

	// Remove binary
	binPath := "/usr/local/bin/vtsbackup"
	if err := os.Remove(binPath); err != nil {
		if os.IsPermission(err) {
			// fmt.Println("Attempting to remove binary with sudo...")
			cmd := exec.Command("sudo", "rm", binPath)
			if out, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to remove backup agent binary: %v\nOutput: %s", err, out)
			}
			// fmt.Println("Binary removed successfully with sudo.")
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove backup agent binary: %v", err)
		}
	} else {
		// fmt.Println("Binary removed successfully.")
	}

	// Remove configuration directory
	configDir := filepath.Join(os.Getenv("HOME"), ".vtsbackup")
	if err := os.RemoveAll(configDir); err != nil {
		if os.IsPermission(err) {
			// fmt.Println("Attempting to remove configuration directory with sudo...")
			cmd := exec.Command("sudo", "rm", "-rf", configDir)
			if out, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to remove backup agent configuration directory: %v\nOutput: %s", err, out)
			}
			// fmt.Println("Configuration directory removed successfully with sudo.")
		} else {
			return fmt.Errorf("failed to remove backup agent configuration directory: %v", err)
		}
	} else {
		// fmt.Println("Configuration directory removed successfully.")
	}

	fmt.Println("Backup agent has been uninstalled successfully.")
	return nil
}

func runBashCommand(command string, args ...string) error {
	// Create a temporary file to store the vts script
	tmpFile, err := os.CreateTemp("", "vts-script")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Read the embedded vts script
	scriptContent, err := vtsScript.ReadFile("vts")
	if err != nil {
		return fmt.Errorf("failed to read embedded vts script: %v", err)
	}

	// Write the script content to the temporary file
	if _, err := tmpFile.Write(scriptContent); err != nil {
		return fmt.Errorf("failed to write vts script to temporary file: %v", err)
	}
	tmpFile.Close()

	// Make the temporary file executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make temporary file executable: %v", err)
	}

	// Execute the script
	cmdArgs := append([]string{tmpFile.Name(), command}, args...)
	cmd := exec.Command("bash", cmdArgs...)

	// Set up pipes for stdin, stdout, and stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	return cmd.Run()
}
