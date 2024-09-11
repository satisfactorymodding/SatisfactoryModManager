package migration

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type migration struct {
	smm2Dir string
}

var Migration *migration

func Init() error {
	if Migration == nil {
		dir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory for migration check: %w", err)
		}
		Migration = &migration{}
		Migration.smm2Dir = dir + "\\SatisfactoryModManager\\profiles\\"
	}
	return nil
}

const migrationSuccessMarkerFile = ".smm3_migration_acknowledged"

// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		slog.Warn("Error when checking path exists, so assuming it does not exist: "+path, slog.Any("error", err))
		return false
	}
}

func (m *migration) NeedsSmm2Migration() bool {
	if pathExists(m.smm2Dir) {
		return !pathExists(m.smm2Dir + migrationSuccessMarkerFile)
	}
	return false
}

func (m *migration) MarkSmm2MigrationSuccess() error {
	file, err := os.Create(m.smm2Dir + migrationSuccessMarkerFile)
	if err != nil {
		return fmt.Errorf("failed to create migration success marker file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}
