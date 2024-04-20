package ficsitcli

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/satisfactorymodding/ficsit-cli/cli/disk"
)

var ServerPicker = &serverPicker{
	disks: make(map[string]disk.Disk),
}

type serverPicker struct {
	disks              map[string]disk.Disk
	nextServerPickerID int
}

type PickerDirectory struct {
	Name           string `json:"name"`
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

func (s *serverPicker) StartPicker(path string) (string, error) {
	id := s.getID()

	d, err := disk.FromPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to create: %w", err)
	}

	s.disks[id] = d

	return id, nil
}

func (s *serverPicker) StopPicker(id string) error {
	if _, ok := s.disks[id]; !ok {
		return fmt.Errorf("no such disk: %s", id)
	}
	delete(s.disks, id)
	return nil
}

func (s *serverPicker) TryPick(id string, path string) (PickerResult, error) {
	d, ok := s.disks[id]
	if !ok {
		return PickerResult{}, fmt.Errorf("no such disk: %s", id)
	}

	result := PickerResult{
		Items: make([]PickerDirectory, 0),
	}

	var err error

	result.IsValidInstall, err = isValidInstall(d, path)
	if err != nil {
		return PickerResult{}, fmt.Errorf("failed to check if valid install: %w", err)
	}

	entries, err := d.ReadDir(path)
	if err != nil {
		return PickerResult{}, fmt.Errorf("failed reading directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		validInstall, err := isValidInstall(d, filepath.Join(path, entry.Name()))
		if err != nil {
			return PickerResult{}, fmt.Errorf("failed to check if valid install: %w", err)
		}

		result.Items = append(result.Items, PickerDirectory{
			Name:           entry.Name(),
			IsValidInstall: validInstall,
		})
	}

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
