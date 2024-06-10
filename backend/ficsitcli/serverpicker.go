package ficsitcli

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/satisfactorymodding/ficsit-cli/cli/disk"
	psUtilDisk "github.com/shirou/gopsutil/v3/disk"
	"golang.org/x/sync/errgroup"
)

type serverPicker struct {
	disks              map[string]diskData
	nextServerPickerID int
}

type diskData struct {
	disk    disk.Disk
	isLocal bool
}

var ServerPicker = &serverPicker{
	disks: make(map[string]diskData),
}

type PickerDirectory struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	IsValidInstall bool   `json:"isValidInstall"`
}

type PickerResult struct {
	IsValidInstall bool              `json:"isValidInstall"`
	Items          []PickerDirectory `json:"items"`
}

func (s *serverPicker) getID() string {
	id := s.nextServerPickerID
	s.nextServerPickerID++
	return strconv.Itoa(id)
}

func (*serverPicker) GetPathSeparator() string {
	return string(filepath.Separator)
}

var remoteSchemes = map[string]struct{}{
	"ftp":  {},
	"sftp": {},
}

func isLocal(path string) (bool, error) {
	parsed, err := url.Parse(path)
	if err != nil {
		return false, fmt.Errorf("failed to parse path: %w", err)
	}

	_, ok := remoteSchemes[parsed.Scheme]
	return !ok, nil
}

func (s *serverPicker) StartPicker(path string) (string, error) {
	id := s.getID()

	d, err := disk.FromPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to create: %w", err)
	}

	local, err := isLocal(path)
	if err != nil {
		return "", fmt.Errorf("failed to check if local: %w", err)
	}

	s.disks[id] = diskData{
		disk:    d,
		isLocal: local,
	}

	return id, nil
}

func (s *serverPicker) StopPicker(id string) error {
	if _, ok := s.disks[id]; !ok {
		return fmt.Errorf("no such disk: %s", id)
	}
	delete(s.disks, id)
	return nil
}

func (s *serverPicker) TryPick(id string, pickPath string) (PickerResult, error) {
	d, ok := s.disks[id]
	if !ok {
		return PickerResult{}, fmt.Errorf("no such disk: %s", id)
	}

	result := PickerResult{
		Items: make([]PickerDirectory, 0),
	}

	if d.isLocal && pickPath == "\\" && runtime.GOOS == "windows" {
		// On windows, the root does not exist, and instead we need to list partitions
		partitions, err := psUtilDisk.Partitions(false)
		if err != nil {
			return PickerResult{}, fmt.Errorf("failed to get partitions: %w", err)
		}
		for _, partition := range partitions {
			validInstall, err := isValidInstall(d.disk, filepath.Join(pickPath, partition.Mountpoint))
			if err != nil {
				return PickerResult{}, fmt.Errorf("failed to check if valid install: %w", err)
			}

			result.Items = append(result.Items, PickerDirectory{
				Name:           partition.Mountpoint,
				Path:           partition.Mountpoint,
				IsValidInstall: validInstall,
			})
		}
		return result, nil
	}

	var err error

	result.IsValidInstall, err = isValidInstall(d.disk, pickPath)
	if err != nil {
		return PickerResult{}, fmt.Errorf("failed to check if valid install: %w", err)
	}

	entries, err := d.disk.ReadDir(pickPath)
	if err != nil {
		return PickerResult{}, fmt.Errorf("failed reading directory: %w", err)
	}

	var wg errgroup.Group
	// We only read the channel after all items have been added,
	// so it must be buffered to avoid a deadlock
	itemsChan := make(chan PickerDirectory, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		wg.TryGo(func() error {
			validInstall, err := isValidInstall(d.disk, filepath.Join(pickPath, entry.Name()))
			if err != nil {
				if errors.Is(err, fs.ErrPermission) {
					return nil
				}
				return fmt.Errorf("failed to check if valid install: %w", err)
			}

			var fullPath string
			if d.isLocal {
				fullPath = filepath.Join(pickPath, entry.Name())
			} else {
				fullPath = path.Join(pickPath, entry.Name())
			}

			itemsChan <- PickerDirectory{
				Name:           entry.Name(),
				Path:           fullPath,
				IsValidInstall: validInstall,
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return PickerResult{}, err //nolint:wrapCheck
	}
	close(itemsChan)

	for item := range itemsChan {
		result.Items = append(result.Items, item)
	}

	slices.SortFunc(result.Items, func(i, j PickerDirectory) int {
		return strings.Compare(i.Name, j.Name)
	})

	return result, nil
}

func isValidInstall(d disk.Disk, path string) (bool, error) {
	var exists bool
	var err error

	exists, err = d.Exists(filepath.Join(path, "FactoryServer.sh"))
	if !exists {
		if err != nil {
			return false, fmt.Errorf("failed reading FactoryServer.sh: %w", err)
		}
	} else {
		return true, nil
	}

	exists, err = d.Exists(filepath.Join(path, "FactoryServer.exe"))
	if !exists {
		if err != nil {
			return false, fmt.Errorf("failed reading FactoryServer.exe: %w", err)
		}
	} else {
		return true, nil
	}

	return false, nil
}
