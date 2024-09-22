package snapshot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hantbk/vtsbackup/logger"
)

type snapshotConfig struct {
	BackupPath  string
	ExcludeList string
	DestDisk    string
	FilePath    string
	FileName    string
	Compression bool
	AutoUnmount bool
}

func Snapshot() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".vtsbackup", "snapshot.conf")
	config, err := loadSnapshotConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load snapshot config: %v", err)
	}

	// Add mount check for destDisk if it's not "/"
	if config.DestDisk != "/" {
		if err := checkAndMountDisk(config.DestDisk); err != nil {
			return fmt.Errorf("failed to check/mount destination disk: %v", err)
		}
	}

	destPath := filepath.Join(config.FilePath, fmt.Sprintf("%s--%s", config.FileName, time.Now().Format("2006-01-02-15-04-05")))

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create a temporary file for excludes
	excludeFile, err := os.CreateTemp("", "vtsbackup-excludes")
	if err != nil {
		return fmt.Errorf("failed to create temporary exclude file: %v", err)
	}
	defer os.Remove(excludeFile.Name())

	if _, err := excludeFile.WriteString(config.ExcludeList); err != nil {
		return fmt.Errorf("failed to write exclude list: %v", err)
	}
	excludeFile.Close()

	// Construct the rsync command
	args := []string{
		"-avAXH",
		"--delete",
		"--info=progress2",
		fmt.Sprintf("--exclude-from=%s", excludeFile.Name()),
		config.BackupPath,
		destPath,
	}

	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Creating snapshot: %s\n", destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsync command failed: %v", err)
	}

	fmt.Printf("Snapshot created successfully at %s\n", destPath)

	if config.Compression {
		if err := compressSnapshot(destPath); err != nil {
			return fmt.Errorf("failed to compress snapshot: %v", err)
		}
	}

	if config.AutoUnmount && config.DestDisk != "/" {
		if err := unmountBackupDisk(config.DestDisk); err != nil {
			return fmt.Errorf("failed to unmount backup disk: %v", err)
		}
		fmt.Printf("Unmounted backup disk: %s\n", config.DestDisk)
	}

	return nil
}

func loadSnapshotConfig(path string) (*snapshotConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	config := &snapshotConfig{
		BackupPath:  "/",
		ExcludeList: "/dev/* /proc/* /sys/* /tmp/* /run/* /mnt/* /media/* /lost+found /cache/ /home/*/.cache/mozilla/* /home/*/.cache/chromium/* /home/*/.local/share/Trash/* /home/*/.gvfs/* *.log",
		DestDisk:    "/mnt/backup",
		FilePath:    "/mnt/backup/snapshot",
		FileName:    "snapshot",
		Compression: false,
		AutoUnmount: false,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "BACKUP_PATH":
			config.BackupPath = value
		case "EXCLUDE_LIST":
			config.ExcludeList = value
		case "DEST_DISK":
			config.DestDisk = value
		case "FILE_PATH":
			config.FilePath = value
		case "FILE_NAME":
			config.FileName = value
		case "COMPRESSION":
			config.Compression, _ = strconv.ParseBool(value)
		case "AUTOUNMOUNT":
			config.AutoUnmount, _ = strconv.ParseBool(value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return config, nil
}

func compressSnapshot(path string) error {
	fmt.Printf("Compressing snapshot: %s\n", path)
	cmd := exec.Command("tar", "czf", path+".tar.gz", "-C", filepath.Dir(path), filepath.Base(path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove uncompressed snapshot: %v", err)
	}
	fmt.Printf("Snapshot compressed: %s.tar.gz\n", path)
	return nil
}

func unmountBackupDisk(disk string) error {
	cmd := exec.Command("umount", disk)
	return cmd.Run()
}

// Add this new function to check and mount the disk
func checkAndMountDisk(destDisk string) error {
	logger := logger.Tag("Snapshot")

	logger.Infof("Checking mount point: %s", destDisk)

	// Check if the disk is mounted
	out, err := exec.Command("mount").Output()
	if err != nil {
		return fmt.Errorf("failed to check mount status: %v", err)
	}

	if !strings.Contains(string(out), destDisk) {
		logger.Info("Disk not mounted, attempting remount")

		cmd := exec.Command("sudo", "mount", "-o", "rw", destDisk)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not mount %s as read-write: %v", destDisk, err)
		}

		// Check again if the mount was successful
		out, err = exec.Command("mount").Output()
		if err != nil || !strings.Contains(string(out), destDisk) {
			return fmt.Errorf("failed to verify mount after attempt: %v", err)
		}
	}

	logger.Info("Disk mounted successfully, continuing...")
	return nil
}
