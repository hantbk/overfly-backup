package snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hantbk/vtsbackup/helper"
	"github.com/spf13/viper"
)

type Config struct {
	SnapshotCount     int
	InfoFile          string
	ShowChangedDirs   string
	RsyncArgs         string
	BackupPermissions string
	LogFile           string
	SpaceErrorLevel   string
	MountDevice       string
	SnapshotRW        string
	Excludes          string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	configFile := helper.AbsolutePath("$HOME/.vtsbackup/snapshot.yml")

	// set config file directly
	if len(configFile) > 0 {
		configFile = helper.AbsolutePath(configFile)
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("snapshot")
		viper.AddConfigPath("$HOME/.vtsbackup")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{
		SnapshotCount:     viper.GetInt("SNAPSHOT_COUNT"),
		InfoFile:          viper.GetString("INFO_FILE"),
		ShowChangedDirs:   viper.GetString("SHOW_CHANGED_DIRS"),
		RsyncArgs:         viper.GetString("RSYNC_ARGS"),
		BackupPermissions: viper.GetString("BACKUP_PERMISSIONS"),
		LogFile:           viper.GetString("LOGFILE"),
		SpaceErrorLevel:   viper.GetString("SPACE_ERRORLEVEL"),
		MountDevice:       viper.GetString("MOUNT_DEVICE"),
		SnapshotRW:        viper.GetString("SNAPSHOT_RW"),
		Excludes:          viper.GetString("EXCLUDES"),
	}

	return config, nil
}

func CreateSnapshot(config *Config, sourcePaths []string) error {
	// Ensure running as root
	if os.Geteuid() != 0 {
		return fmt.Errorf("sorry, must be root")
	}

	// Ensure destination path exists
	if _, err := os.Stat(config.SnapshotRW); os.IsNotExist(err) {
		return fmt.Errorf("destination path %s not found", config.SnapshotRW)
	}

	// Remount the RW mount point as RW
	if err := runCommand("mount", "-o", "remount,rw", config.MountDevice, config.SnapshotRW); err != nil {
		return fmt.Errorf("snapshot: could not remount %s readwrite", config.SnapshotRW)
	}

	// Create snapshot directories if needed
	for i := 0; i < config.SnapshotCount; i++ {
		snapshotDir := fmt.Sprintf("%s/snapshot.%d", config.SnapshotRW, i)
		if _, err := os.Stat(snapshotDir); os.IsNotExist(err) {
			if err := os.Mkdir(snapshotDir, 0750); err != nil {
				return fmt.Errorf("failed to create snapshot directory %s: %w", snapshotDir, err)
			}
		}
	}

	// Delete oldest snapshot
	oldestSnapshot := fmt.Sprintf("%s/snapshot.%d", config.SnapshotRW, config.SnapshotCount-1)
	if err := os.RemoveAll(oldestSnapshot); err != nil {
		return fmt.Errorf("failed to delete oldest snapshot %s: %w", oldestSnapshot, err)
	}

	// Renumber snapshots
	for i := config.SnapshotCount - 1; i > 0; i-- {
		oldDir := fmt.Sprintf("%s/snapshot.%d", config.SnapshotRW, i-1)
		newDir := fmt.Sprintf("%s/snapshot.%d", config.SnapshotRW, i)
		if err := os.Rename(oldDir, newDir); err != nil {
			return fmt.Errorf("failed to rename snapshot directory from %s to %s: %w", oldDir, newDir, err)
		}
	}

	// Make a hard-link-only copy of the latest snapshot
	latestSnapshot := fmt.Sprintf("%s/snapshot.0", config.SnapshotRW)
	if _, err := os.Stat(latestSnapshot); err == nil {
		if err := runCommand("cp", "-al", latestSnapshot, fmt.Sprintf("%s/snapshot.1", config.SnapshotRW)); err != nil {
			return fmt.Errorf("failed to create hard-link-only copy of the latest snapshot: %w", err)
		}
	}

	// Perform rsync to create new snapshot
	rsyncArgs := strings.Fields(config.RsyncArgs)
	rsyncArgs = append(rsyncArgs, "--delete", "--delete-excluded", fmt.Sprintf("--exclude-from=%s", config.Excludes))
	rsyncArgs = append(rsyncArgs, sourcePaths...)
	rsyncArgs = append(rsyncArgs, fmt.Sprintf("%s/snapshot.0/", config.SnapshotRW))

	cmd := exec.Command("rsync", rsyncArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("rsync command failed: %s, output: %s", err, output)
	}

	// Log completion
	logMessage := fmt.Sprintf("Backup completed at %s", time.Now().Format(time.RFC3339))
	if err := os.WriteFile(fmt.Sprintf("%s/snapshot.0/%s", config.SnapshotRW, config.InfoFile), []byte(logMessage), 0644); err != nil {
		return fmt.Errorf("failed to write info file: %w", err)
	}

	// Remount the RW snapshot mount point as readonly
	if err := runCommand("mount", "-o", "remount,ro", config.MountDevice, config.SnapshotRW); err != nil {
		return fmt.Errorf("snapshot: could not remount %s readonly", config.SnapshotRW)
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
