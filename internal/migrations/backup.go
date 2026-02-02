package migrations

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupResult contains information about a created backup.
type BackupResult struct {
	// BackupDir is the full path to the backup directory.
	BackupDir string
	// Timestamp is when the backup was created.
	Timestamp time.Time
	// NodeCount is the number of node files backed up.
	NodeCount int
}

// CreateBackup creates a backup of the current project state.
// It creates a .deco/backup-<timestamp>/ directory containing:
// - config.yaml (copy of current config)
// - nodes/ (copy of all node files)
func CreateBackup(rootDir string) (*BackupResult, error) {
	timestamp := time.Now()
	backupName := fmt.Sprintf("backup-%s", timestamp.Format("20060102-150405"))
	backupDir := filepath.Join(rootDir, ".deco", backupName)

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy config.yaml
	configSrc := filepath.Join(rootDir, ".deco", "config.yaml")
	configDst := filepath.Join(backupDir, "config.yaml")
	if err := copyFile(configSrc, configDst); err != nil {
		return nil, fmt.Errorf("failed to backup config: %w", err)
	}

	// Copy nodes directory
	nodesSrc := filepath.Join(rootDir, ".deco", "nodes")
	nodesDst := filepath.Join(backupDir, "nodes")
	nodeCount, err := copyDir(nodesSrc, nodesDst)
	if err != nil {
		return nil, fmt.Errorf("failed to backup nodes: %w", err)
	}

	return &BackupResult{
		BackupDir: backupDir,
		Timestamp: timestamp,
		NodeCount: nodeCount,
	}, nil
}

// RestoreBackup restores a project from a backup directory.
// It replaces the current config.yaml and nodes/ with the backup contents.
func RestoreBackup(rootDir, backupDir string) error {
	// Verify backup exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return fmt.Errorf("backup directory not found: %s", backupDir)
	}

	// Verify backup has required files
	backupConfig := filepath.Join(backupDir, "config.yaml")
	if _, err := os.Stat(backupConfig); os.IsNotExist(err) {
		return fmt.Errorf("backup config not found: %s", backupConfig)
	}

	// Restore config.yaml
	configDst := filepath.Join(rootDir, ".deco", "config.yaml")
	if err := copyFile(backupConfig, configDst); err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	// Restore nodes directory
	backupNodes := filepath.Join(backupDir, "nodes")
	nodesDst := filepath.Join(rootDir, ".deco", "nodes")

	// Remove current nodes directory
	if err := os.RemoveAll(nodesDst); err != nil {
		return fmt.Errorf("failed to remove current nodes: %w", err)
	}

	// Copy backup nodes if they exist
	if _, err := os.Stat(backupNodes); err == nil {
		if _, err := copyDir(backupNodes, nodesDst); err != nil {
			return fmt.Errorf("failed to restore nodes: %w", err)
		}
	}

	return nil
}

// ListBackups returns a list of available backup directories.
func ListBackups(rootDir string) ([]string, error) {
	decoDir := filepath.Join(rootDir, ".deco")
	entries, err := os.ReadDir(decoDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 7 && entry.Name()[:7] == "backup-" {
			backups = append(backups, filepath.Join(decoDir, entry.Name()))
		}
	}
	return backups, nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	// Check if source exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Source doesn't exist, nothing to copy
		}
		return err
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Preserve permissions
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir recursively copies a directory from src to dst.
// Returns the number of files copied.
func copyDir(src, dst string) (int, error) {
	// Check if source exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil // Source doesn't exist, nothing to copy
		}
		return 0, err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return 0, err
	}

	var count int

	// Walk source directory
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		if err := copyFile(path, dstPath); err != nil {
			return err
		}
		count++
		return nil
	})

	return count, err
}
