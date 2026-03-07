//go:build windows
// +build windows

package service

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

func (s *StorageService) checkDiskSpace() error {
	uploadPath := s.uploadPath

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		uploadPath = "."
	}

	pathPtr, err := syscall.UTF16PtrFromString(uploadPath)
	if err != nil {
		return nil
	}

	var freeBytesAvailable uint64
	var totalBytes uint64
	var freeBytes uint64

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&freeBytes)),
	)

	if ret == 0 {
		return nil
	}

	usedBytes := totalBytes - freeBytesAvailable
	usagePercent := int((usedBytes * 100) / totalBytes)

	if usagePercent >= s.maxDiskUsage {
		return errors.New("disk usage exceeds maximum allowed percentage")
	}

	return nil
}
