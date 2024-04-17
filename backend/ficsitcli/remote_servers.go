package ficsitcli

import (
	"fmt"
	"log/slog"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

func (f *ficsitCLI) GetRemoteInstallations() []string {
	paths := make([]string, 0, f.installationMetadata.Size())
	for _, install := range f.GetInstallations() {
		meta, ok := f.installationMetadata.Load(install)
		if ok && meta.Info != nil {
			if meta.Info.Location != common.LocationTypeRemote {
				continue
			}
		}
		// Missing metadata means an unavailable remote installation
		paths = append(paths, install)
	}
	return paths
}

func (f *ficsitCLI) AddRemoteServer(path string) error {
	if f.ficsitCli.Installations.GetInstallation(path) != nil {
		return fmt.Errorf("installation already exists")
	}
	l := slog.With(slog.String("task", "addRemoteServer"), slog.String("path", path))

	installation, err := f.ficsitCli.Installations.AddInstallation(f.ficsitCli, path, f.GetFallbackProfile())
	if err != nil {
		return fmt.Errorf("failed to add installation: %w", err)
	}

	err = f.ficsitCli.Installations.Save()
	if err != nil {
		l.Error("failed to save installations", slog.Any("error", err))
	}

	meta, err := f.getRemoteServerMetadata(installation)
	if err != nil {
		return fmt.Errorf("failed to get remote server metadata: %w", err)
	}

	f.installationMetadata.Store(path, installationMetadata{
		State: InstallStateValid,
		Info:  meta,
	})

	f.EmitGlobals()

	return nil
}

func (f *ficsitCLI) RemoveRemoteServer(path string) error {
	metadata, ok := f.installationMetadata.Load(path)
	if !ok {
		return fmt.Errorf("installation not found")
	}
	if metadata.State == InstallStateLoading {
		return fmt.Errorf("installation is still loading")
	}
	if metadata.Info != nil && metadata.Info.Location != common.LocationTypeRemote {
		return fmt.Errorf("installation is not remote")
	}
	err := f.ficsitCli.Installations.DeleteInstallation(path)
	if err != nil {
		return fmt.Errorf("failed to delete installation: %w", err)
	}
	err = f.ficsitCli.Installations.Save()
	if err != nil {
		slog.Error("failed to save installations", slog.Any("error", err))
	}
	f.installationMetadata.Delete(path)
	f.EmitGlobals()
	return nil
}
