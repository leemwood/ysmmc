//go:build !windows
// +build !windows

package service

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

func (s *StorageService) checkDiskSpace() error {
	uploadPath := s.uploadPath

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		uploadPath = "."
	}

	var stat unix.Statfs_t
	if err := unix.Statfs(uploadPath, &stat); err != nil {
		return nil
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes
	usagePercent := int((usedBytes * 100) / totalBytes)

	if usagePercent >= s.maxDiskUsage {
		return errors.New("disk usage exceeds maximum allowed percentage")
	}

	return nil
}
